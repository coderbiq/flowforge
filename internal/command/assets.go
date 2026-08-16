package command

import (
	"fmt"
	"path/filepath"

	"flowforge/internal/core"
	"flowforge/internal/version"
)

func applyAssetUpdates(projectRoot string, adopt bool) (*AssetUpdateReport, error) {
	oldManifest, err := core.LoadProjectManifest(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("loading existing project manifest: %w", err)
	}

	newManifest, err := core.GenerateManifest(embeddedAssets, version.Version)
	if err != nil {
		return nil, fmt.Errorf("generating manifest: %w", err)
	}
	for _, entry := range oldManifest.Files {
		if entry.Type == "opencode_agent" || entry.Type == "codex_agent" || entry.Type == "orchestration_block" {
			newManifest.Files = append(newManifest.Files, entry)
		}
	}
	newManifest.Mode = oldManifest.Mode
	newManifest.HostIntent = oldManifest.HostIntent
	newManifest.Renderer = oldManifest.Renderer
	newManifest.DisabledHosts = append([]string(nil), oldManifest.DisabledHosts...)
	newManifest.PendingHosts = append([]string(nil), oldManifest.PendingHosts...)
	// Dynamic entries are owned by the host-aware sync lifecycle. Keep every
	// registered entry, including dormant and disabled ones, in the generated
	// manifest so an asset upgrade cannot enable, delete, or re-baseline it.

	diff := core.CompareManifests(oldManifest, newManifest, projectRoot)
	if adopt && len(diff.Conflict) > 0 {
		oldBySource := make(map[string]core.FileEntry, len(oldManifest.Files))
		for _, entry := range oldManifest.Files {
			oldBySource[entry.Source] = entry
		}
		adopted := make([]core.FileEntry, 0, len(diff.Conflict))
		preserved := make([]core.FileEntry, 0, len(diff.Conflict))
		for _, entry := range diff.Conflict {
			// --adopt is bounded to files already identified by the trusted
			// manifest. A newly appearing arbitrary target remains preserved.
			if _, known := oldBySource[entry.Source]; known {
				adopted = append(adopted, entry)
			} else {
				preserved = append(preserved, entry)
			}
		}
		diff.Updated = append(diff.Updated, adopted...)
		diff.Conflict = preserved
	}
	if !diff.HasChanges() {
		return nil, nil
	}

	backupDir := ""
	if oldManifest.CLIVersion != "" {
		backupDir = filepath.Join(projectRoot, ".flowforge", "backup", oldManifest.CLIVersion)
	}
	report := core.ApplyUpgrade(diff, newManifest, projectRoot, embeddedAssets, backupDir)

	if report.Error != nil {
		return nil, report.Error
	}
	if len(diff.Conflict) > 0 {
		oldBySource := make(map[string]core.FileEntry, len(oldManifest.Files))
		for _, entry := range oldManifest.Files {
			oldBySource[entry.Source] = entry
		}
		conflicts := make(map[string]bool, len(diff.Conflict))
		for _, entry := range diff.Conflict {
			conflicts[entry.Source] = true
		}
		kept := make([]core.FileEntry, 0, len(newManifest.Files))
		for _, entry := range newManifest.Files {
			if !conflicts[entry.Source] {
				kept = append(kept, entry)
				continue
			}
			if old, ok := oldBySource[entry.Source]; ok {
				kept = append(kept, old)
			}
		}
		newManifest.Files = kept
	}

	if err := newManifest.Save(projectRoot); err != nil {
		return nil, fmt.Errorf("saving updated manifest: %w", err)
	}

	return &AssetUpdateReport{
		SummaryLine:  diff.Summary(),
		BlockUpdated: report.BlockUpdated,
		Added:        report.Added,
		Updated:      report.Updated,
		Conflict:     report.Conflict,
	}, nil
}

type AssetUpdateReport struct {
	SummaryLine  string
	BlockUpdated bool
	Added        []core.FileEntry
	Updated      []core.FileEntry
	Conflict     []core.FileEntry
}

func (r *AssetUpdateReport) Summary() string { return r.SummaryLine }
