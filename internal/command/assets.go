package command

import (
	"fmt"
	"path/filepath"

	"flowforge/internal/core"
	"flowforge/internal/version"
)

func applyAssetUpdates(projectRoot string) (*AssetUpdateReport, error) {
	oldManifest, err := core.LoadProjectManifest(projectRoot)
	if err != nil {
		oldManifest = &core.ProjectManifest{}
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
	newManifest.DisabledHosts = append([]string(nil), oldManifest.DisabledHosts...)
	newManifest.PendingHosts = append([]string(nil), oldManifest.PendingHosts...)

	diff := core.CompareManifests(oldManifest, newManifest, projectRoot)
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
		for i, entry := range newManifest.Files {
			if conflicts[entry.Source] {
				newManifest.Files[i] = oldBySource[entry.Source]
			}
		}
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
