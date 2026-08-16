package uninstall

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"flowforge/internal/core"
)

type Result struct {
	Removed []string
	Errors  []error
}

// Project filesystem operations are variables so lifecycle tests can exercise
// every failure boundary without relying on host permissions.
var (
	projectMkdirAll     = os.MkdirAll
	projectRemove       = os.Remove
	projectWriteFile    = os.WriteFile
	projectSaveManifest = core.SaveManifestAtomic
)

func CleanBinary() (*Result, error) {
	result := &Result{}

	exe, err := os.Executable()
	if err != nil {
		return result, fmt.Errorf("finding executable: %w", err)
	}

	real, err := filepath.EvalSymlinks(exe)
	if err != nil {
		real = exe
	}

	if err := os.Remove(real); err != nil {
		if !os.IsNotExist(err) {
			result.Errors = append(result.Errors, fmt.Errorf("removing binary %s: %w", real, err))
		}
	} else {
		result.Removed = append(result.Removed, real)
	}

	oldPath := real + ".old"
	if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
		result.Errors = append(result.Errors, fmt.Errorf("removing backup %s: %w", oldPath, err))
	}

	return result, nil
}

func CleanConfig(homeDir string) (*Result, error) {
	result := &Result{}

	configDir := filepath.Join(homeDir, ".flowforge")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return result, nil
	}

	if err := os.RemoveAll(configDir); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("removing config dir: %w", err))
	} else {
		result.Removed = append(result.Removed, configDir)
	}

	return result, nil
}

func CleanProject(projectRoot string) (*Result, error) {
	result := &Result{}
	loaded, err := core.LoadProjectManifestMigration(projectRoot)
	if err != nil {
		return result, fmt.Errorf("loading project manifest for uninstall: %w", err)
	}
	if loaded.Manifest.Version != core.ManifestVersionV2 {
		return result, fmt.Errorf("uninstall requires a validated v2 manifest")
	}

	entries := append([]core.FileEntry(nil), loaded.Manifest.Files...)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Target < entries[j].Target })
	dynamic := make([]projectCleanupEntry, 0)
	static := make([]projectCleanupEntry, 0)
	for _, entry := range entries {
		path, pathErr := core.ProjectPath(projectRoot, entry.Target)
		if pathErr != nil {
			return result, fmt.Errorf("validating uninstall target %q: %w", entry.Target, pathErr)
		}
		item := projectCleanupEntry{entry: entry, path: path}
		if isDynamicEntry(entry) {
			dynamic = append(dynamic, item)
		} else {
			static = append(static, item)
		}
	}

	// The base FLOWFORGE block is a static asset and is deliberately retained.
	// Only the separately registered orchestration marker is removable.
	for i := range static {
		item := &static[i]
		if item.entry.Type == "agents_block" {
			continue
		}
		data, readErr := os.ReadFile(item.path)
		if os.IsNotExist(readErr) {
			continue
		}
		if readErr != nil {
			return result, fmt.Errorf("reading static uninstall target %q: %w", item.entry.Target, readErr)
		}
		if core.SHA256Hex(data) != item.entry.SHA256 {
			result.Errors = append(result.Errors, fmt.Errorf("preserving modified static target %q", item.entry.Target))
			continue
		}
		item.data = data
	}

	backupDir, err := allocateBackupDir(projectRoot)
	if err != nil {
		return result, err
	}
	backups := make([]projectCleanupEntry, 0, len(dynamic))
	for _, item := range dynamic {
		data, readErr := os.ReadFile(item.path)
		if os.IsNotExist(readErr) {
			continue
		}
		if readErr != nil {
			return result, fmt.Errorf("reading dynamic uninstall target %q: %w", item.entry.Target, readErr)
		}
		item.data = data
		backups = append(backups, item)
	}
	if err := backupEntries(backupDir, backups); err != nil {
		return result, err
	}

	removed := make([]projectCleanupEntry, 0, len(static)+len(dynamic))
	rollback := func() error {
		for i := len(removed) - 1; i >= 0; i-- {
			item := removed[i]
			if err := os.WriteFile(item.path, item.data, 0644); err != nil {
				return fmt.Errorf("restoring %q: %w", item.entry.Target, err)
			}
		}
		return nil
	}
	for _, item := range static {
		if item.entry.Type == "agents_block" || item.data == nil {
			continue
		}
		if err := projectRemove(item.path); err != nil {
			if rollbackErr := rollback(); rollbackErr != nil {
				return result, fmt.Errorf("removing static target %q: %w (rollback failed: %v)", item.entry.Target, err, rollbackErr)
			}
			return result, fmt.Errorf("removing static target %q: %w", item.entry.Target, err)
		}
		removed = append(removed, item)
		result.Removed = append(result.Removed, item.entry.Target)
	}
	for _, item := range backups {
		updated := item.data
		if item.entry.Type == "orchestration_block" {
			if item.entry.Markers == nil {
				return result, fmt.Errorf("orchestration entry %q has no markers", item.entry.Target)
			}
			var found bool
			updated, found, err = core.RemoveMarkedBlockContent(updated, item.entry.Markers.Start, item.entry.Markers.End)
			if err != nil {
				if rollbackErr := rollback(); rollbackErr != nil {
					return result, fmt.Errorf("removing orchestration block from %q: %w (rollback failed: %v)", item.entry.Target, err, rollbackErr)
				}
				return result, fmt.Errorf("removing orchestration block from %q: %w", item.entry.Target, err)
			}
			if !found {
				result.Errors = append(result.Errors, fmt.Errorf("preserving modified orchestration target %q", item.entry.Target))
				continue
			}
			if err := projectWriteFile(item.path, updated, 0644); err != nil {
				if rollbackErr := rollback(); rollbackErr != nil {
					return result, fmt.Errorf("updating orchestration target %q: %w (rollback failed: %v)", item.entry.Target, err, rollbackErr)
				}
				return result, fmt.Errorf("updating orchestration target %q: %w", item.entry.Target, err)
			}
		} else if err := projectRemove(item.path); err != nil {
			if rollbackErr := rollback(); rollbackErr != nil {
				return result, fmt.Errorf("removing dynamic target %q: %w (rollback failed: %v)", item.entry.Target, err, rollbackErr)
			}
			return result, fmt.Errorf("removing dynamic target %q: %w", item.entry.Target, err)
		}
		removed = append(removed, projectCleanupEntry{entry: item.entry, path: item.path, data: item.data})
		result.Removed = append(result.Removed, item.entry.Target)
	}

	kept := make([]core.FileEntry, 0, 1)
	for _, entry := range loaded.Manifest.Files {
		if entry.Type == "agents_block" {
			kept = append(kept, entry)
		}
	}
	loaded.Manifest.Files = kept
	if err := projectSaveManifest(projectRoot, loaded.Manifest); err != nil {
		if rollbackErr := rollback(); rollbackErr != nil {
			return result, fmt.Errorf("saving uninstall manifest: %w (rollback failed: %v)", err, rollbackErr)
		}
		return result, fmt.Errorf("saving uninstall manifest: %w", err)
	}
	return result, nil
}

type projectCleanupEntry struct {
	entry core.FileEntry
	path  string
	data  []byte
}

func isDynamicEntry(entry core.FileEntry) bool {
	return entry.Type == "opencode_agent" || entry.Type == "codex_agent" || entry.Type == "orchestration_block"
}

func allocateBackupDir(projectRoot string) (string, error) {
	base := filepath.Join(projectRoot, ".flowforge", "backups", "uninstall")
	stamp := time.Now().UTC().Format("20060102T150405Z")
	for n := 0; ; n++ {
		name := stamp
		if n > 0 {
			name = fmt.Sprintf("%s-%d", stamp, n)
		}
		candidate := filepath.Join(base, name)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate, nil
		} else if err != nil {
			return "", fmt.Errorf("checking uninstall backup directory: %w", err)
		}
	}
}

func backupEntries(dir string, entries []projectCleanupEntry) error {
	if len(entries) == 0 {
		return nil
	}
	for _, item := range entries {
		if err := projectMkdirAll(filepath.Dir(filepath.Join(dir, filepath.FromSlash(item.entry.Target))), 0755); err != nil {
			return fmt.Errorf("creating uninstall backup directory: %w", err)
		}
		path := filepath.Join(dir, filepath.FromSlash(item.entry.Target))
		if err := projectWriteFile(path, item.data, 0644); err != nil {
			return fmt.Errorf("backing up %q: %w", item.entry.Target, err)
		}
	}
	return nil
}

func isEmptyDir(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) == 0
}

func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *Result) ErrorSummary() string {
	if len(r.Errors) == 0 {
		return ""
	}
	var msg string
	for _, e := range r.Errors {
		msg += "  " + e.Error() + "\n"
	}
	return msg
}
