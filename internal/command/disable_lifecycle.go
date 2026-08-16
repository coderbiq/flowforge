package command

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"flowforge/internal/core"
)

type BackupEntry struct {
	Target      string
	Destination string
	Status      string
}

type DisablePlan struct {
	Hosts               []string
	Entries             []core.FileEntry
	Statuses            []string
	BackupDir           string
	Backups             []BackupEntry
	RemoveOrchestration bool
}

var (
	disableWriteFile = os.WriteFile
	disableRemove    = os.Remove
	disableMkdirAll  = os.MkdirAll
	disableNow       = time.Now
)

func buildDisablePlan(root string, manifest *core.ProjectManifest, requested []string, dryRun bool) (*DisablePlan, error) {
	if manifest == nil || manifest.Version != core.ManifestVersionV2 {
		return nil, fmt.Errorf("disable requires a validated v2 manifest")
	}
	hosts := makeHostSet(requested)
	if len(hosts) == 0 {
		hosts[core.HostOpenCode] = true
		hosts[core.HostCodex] = true
	}
	plan := &DisablePlan{Hosts: sortedHosts(hosts)}
	remainingOpenCode := manifest.HostIntent.OpenCode == core.HostEnabled && !hosts[core.HostOpenCode]
	remainingCodex := manifest.HostIntent.Codex == core.HostEnabled && !hosts[core.HostCodex]
	plan.RemoveOrchestration = !remainingOpenCode && !remainingCodex
	for _, entry := range manifest.Files {
		if !isDisableEntry(entry, hosts, plan.RemoveOrchestration) {
			continue
		}
		if entry.Target == "AGENTS.md" && entry.Type == "orchestration_block" {
			// The orchestration block is shared by enabled hosts and is removed
			// as a marker block; it is never treated as a whole-file deletion.
			plan.Entries = append(plan.Entries, entry)
			continue
		}
		plan.Entries = append(plan.Entries, entry)
	}
	sort.Slice(plan.Entries, func(i, j int) bool { return plan.Entries[i].Target < plan.Entries[j].Target })
	for _, entry := range plan.Entries {
		path, err := core.ProjectPath(root, entry.Target)
		if err != nil {
			return nil, fmt.Errorf("validating disable target %q: %w", entry.Target, err)
		}
		data, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			plan.Statuses = append(plan.Statuses, entry.Target+": missing/already absent")
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("reading disable target %q: %w", entry.Target, err)
		}
		status := "clean"
		if entry.Markers != nil {
			block, found, markerErr := core.ExtractMarkedBlock(data, entry.Markers.Start, entry.Markers.End)
			if markerErr != nil {
				return nil, fmt.Errorf("checking markers for %q: %w", entry.Target, markerErr)
			}
			if !found {
				status = "modified-but-authorized"
			} else if core.SHA256Hex(block) != entry.SHA256 {
				status = "modified-but-authorized"
			}
		} else if core.SHA256Hex(data) != entry.SHA256 {
			status = "modified-but-authorized"
		}
		plan.Statuses = append(plan.Statuses, entry.Target+": "+status)
		if !dryRun {
			plan.Backups = append(plan.Backups, BackupEntry{Target: entry.Target, Status: status})
		}
	}
	if dryRun {
		plan.BackupDir = filepath.Join(root, ".flowforge", "backups", "subagent-disable", "<timestamp>")
		return plan, nil
	}
	dir, err := allocateDisableBackupDir(root)
	if err != nil {
		return nil, err
	}
	plan.BackupDir = dir
	for i := range plan.Backups {
		plan.Backups[i].Destination = filepath.Join(dir, filepath.FromSlash(plan.Backups[i].Target))
	}
	return plan, nil
}

func isDisableEntry(entry core.FileEntry, hosts hostSet, removeOrchestration bool) bool {
	if entry.Type != "opencode_agent" && entry.Type != "codex_agent" && entry.Type != "orchestration_block" {
		return false
	}
	if entry.Type == "opencode_agent" {
		return hosts[core.HostOpenCode]
	}
	if entry.Type == "codex_agent" {
		return hosts[core.HostCodex]
	}
	return removeOrchestration
}

func allocateDisableBackupDir(root string) (string, error) {
	base := filepath.Join(root, ".flowforge", "backups", "subagent-disable")
	stamp := disableNow().UTC().Format("20060102T150405Z")
	for n := 0; ; n++ {
		name := stamp
		if n > 0 {
			name = fmt.Sprintf("%s-%d", stamp, n)
		}
		candidate := filepath.Join(base, name)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate, nil
		} else if err != nil {
			return "", fmt.Errorf("checking backup directory: %w", err)
		}
	}
}

func backupDisablePlan(root string, plan *DisablePlan) error {
	if len(plan.Backups) == 0 {
		return nil
	}
	cleanupOnError := func(cause error) error {
		if cleanupErr := os.RemoveAll(plan.BackupDir); cleanupErr != nil {
			return fmt.Errorf("%w (cleaning backup directory failed: %v)", cause, cleanupErr)
		}
		removeIfEmpty := func(dir string) error {
			entries, err := os.ReadDir(dir)
			if os.IsNotExist(err) || (err == nil && len(entries) > 0) {
				return nil
			}
			if err != nil {
				return err
			}
			if err := os.Remove(dir); err != nil && !os.IsNotExist(err) {
				return err
			}
			return nil
		}
		backupRoot := filepath.Dir(filepath.Dir(plan.BackupDir))
		if cleanupErr := removeIfEmpty(filepath.Dir(plan.BackupDir)); cleanupErr != nil {
			return fmt.Errorf("%w (cleaning empty disable backup parent failed: %v)", cause, cleanupErr)
		}
		if cleanupErr := removeIfEmpty(backupRoot); cleanupErr != nil {
			return fmt.Errorf("%w (cleaning empty backup parent failed: %v)", cause, cleanupErr)
		}
		return cause
	}
	if err := disableMkdirAll(plan.BackupDir, 0755); err != nil {
		return fmt.Errorf("creating disable backup directory: %w", err)
	}
	for i := range plan.Backups {
		entry := &plan.Backups[i]
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(entry.Target)))
		if err != nil {
			return cleanupOnError(fmt.Errorf("reading backup source %q: %w", entry.Target, err))
		}
		if err := disableMkdirAll(filepath.Dir(entry.Destination), 0755); err != nil {
			return cleanupOnError(fmt.Errorf("creating backup parent for %q: %w", entry.Target, err))
		}
		if err := disableWriteFile(entry.Destination, data, 0644); err != nil {
			return cleanupOnError(fmt.Errorf("backing up %q: %w", entry.Target, err))
		}
	}
	return nil
}

func executeDisable(root string, manifest *core.ProjectManifest, plan *DisablePlan) error {
	if err := backupDisablePlan(root, plan); err != nil {
		return err
	}
	manifestPath := filepath.Join(root, ".flowforge", core.ManifestFileName)
	oldManifest, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("reading manifest snapshot: %w", err)
	}
	removed := make([]BackupEntry, 0, len(plan.Backups))
	rollback := func() error {
		for i := len(removed) - 1; i >= 0; i-- {
			item := removed[i]
			data, readErr := os.ReadFile(item.Destination)
			if readErr != nil {
				return fmt.Errorf("reading rollback backup %q: %w", item.Target, readErr)
			}
			if writeErr := disableWriteFile(filepath.Join(root, filepath.FromSlash(item.Target)), data, 0644); writeErr != nil {
				return fmt.Errorf("restoring %q: %w", item.Target, writeErr)
			}
		}
		if writeErr := disableWriteFile(manifestPath, oldManifest, 0644); writeErr != nil {
			return fmt.Errorf("restoring manifest: %w", writeErr)
		}
		return nil
	}
	for _, entry := range plan.Entries {
		path := filepath.Join(root, filepath.FromSlash(entry.Target))
		if entry.Type == "orchestration_block" {
			if entry.Markers == nil {
				return fmt.Errorf("orchestration entry %q has no markers", entry.Target)
			}
			data, readErr := os.ReadFile(path)
			if os.IsNotExist(readErr) {
				continue
			}
			if readErr != nil {
				return readErr
			}
			updated, found, blockErr := core.RemoveMarkedBlockContent(data, entry.Markers.Start, entry.Markers.End)
			if blockErr != nil {
				return blockErr
			}
			if found {
				if writeErr := disableWriteFile(path, updated, 0644); writeErr != nil {
					if rollbackErr := rollback(); rollbackErr != nil {
						return fmt.Errorf("removing orchestration block: %w (rollback failed: %v)", writeErr, rollbackErr)
					}
					return writeErr
				}
				removed = append(removed, BackupEntry{Target: entry.Target, Destination: filepath.Join(plan.BackupDir, filepath.FromSlash(entry.Target))})
			}
			continue
		}
		_, statErr := os.Stat(path)
		if os.IsNotExist(statErr) {
			continue
		}
		if statErr != nil {
			if rollbackErr := rollback(); rollbackErr != nil {
				return fmt.Errorf("checking %q: %w (rollback failed: %v)", entry.Target, statErr, rollbackErr)
			}
			return statErr
		}
		if removeErr := disableRemove(path); removeErr != nil {
			if rollbackErr := rollback(); rollbackErr != nil {
				return fmt.Errorf("deleting %q: %w (rollback failed: %v)", entry.Target, removeErr, rollbackErr)
			}
			return fmt.Errorf("deleting %q: %w", entry.Target, removeErr)
		}
		removed = append(removed, BackupEntry{Target: entry.Target, Destination: filepath.Join(plan.BackupDir, filepath.FromSlash(entry.Target))})
	}
	kept := make([]core.FileEntry, 0, len(manifest.Files))
	for _, entry := range manifest.Files {
		if isDisableEntry(entry, makeHostSet(plan.Hosts), plan.RemoveOrchestration) {
			continue
		}
		kept = append(kept, entry)
	}
	manifest.Files = kept
	if contains(plan.Hosts, core.HostOpenCode) {
		manifest.HostIntent.OpenCode = core.HostDisabled
	}
	if contains(plan.Hosts, core.HostCodex) {
		manifest.HostIntent.Codex = core.HostDisabled
	}
	if manifest.HostIntent.OpenCode != core.HostEnabled && manifest.HostIntent.Codex != core.HostEnabled {
		manifest.Mode = core.ManifestModeNonSubagent
	}
	manifest.DisabledHosts = disabledHostList(manifest)
	if err := saveManifest(root, manifest); err != nil {
		if rollbackErr := rollback(); rollbackErr != nil {
			return fmt.Errorf("saving manifest: %w (rollback failed: %v)", err, rollbackErr)
		}
		return fmt.Errorf("saving manifest: %w", err)
	}
	return nil
}

func contains(values []string, value string) bool {
	return strings.Contains("\n"+strings.Join(values, "\n")+"\n", "\n"+value+"\n")
}

func disableProject(root string, hosts []string, dryRun bool) (*DisablePlan, error) {
	result, err := core.LoadProjectManifestMigration(root)
	if err != nil {
		return nil, fmt.Errorf("loading project manifest for disable: %w", err)
	}
	manifest := result.Manifest
	plan, err := buildDisablePlan(root, manifest, hosts, dryRun)
	if err != nil {
		return nil, err
	}
	if dryRun {
		return plan, nil
	}
	if err := executeDisable(root, manifest, plan); err != nil {
		return nil, err
	}
	return plan, nil
}
