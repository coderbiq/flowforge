package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type UpgradeReport struct {
	Added    []FileEntry
	Updated  []FileEntry
	Conflict []FileEntry
	BlockUpdated bool
	Error    error
}

func ApplyUpgrade(diff *DiffResult, newManifest *ProjectManifest, projectRoot string, assetsFS fs.FS, backupDir string) *UpgradeReport {
	report := &UpgradeReport{}

	if backupDir != "" {
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			report.Error = fmt.Errorf("creating backup dir: %w", err)
			return report
		}
	}

	for _, entry := range diff.Added {
		if entry.Type == "agents_block" {
			if err := applyAgentsBlockEntry(entry, assetsFS, backupDir); err != nil {
				report.Error = fmt.Errorf("applying agents block: %w", err)
				return report
			}
			report.BlockUpdated = true
			continue
		}

		if err := applyAddedFile(entry, projectRoot, assetsFS, backupDir); err != nil {
			report.Error = fmt.Errorf("adding file %s: %w", entry.Target, err)
			return report
		}
		report.Added = append(report.Added, entry)
	}

	for _, entry := range diff.Updated {
		if entry.Type == "agents_block" {
			if err := applyAgentsBlockEntry(entry, assetsFS, backupDir); err != nil {
				report.Error = fmt.Errorf("updating agents block: %w", err)
				return report
			}
			report.BlockUpdated = true
			continue
		}

		if err := applyUpdatedFile(entry, projectRoot, assetsFS, backupDir); err != nil {
			report.Error = fmt.Errorf("updating file %s: %w", entry.Target, err)
			return report
		}
		report.Updated = append(report.Updated, entry)
	}

	report.Conflict = append(report.Conflict, diff.Conflict...)

	return report
}

func applyAddedFile(entry FileEntry, projectRoot string, assetsFS fs.FS, backupDir string) error {
	content, err := fs.ReadFile(assetsFS, entry.Source)
	if err != nil {
		return fmt.Errorf("reading asset %s: %w", entry.Source, err)
	}

	targetPath := filepath.Join(projectRoot, entry.Target)
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating target dir: %w", err)
	}

	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

func applyUpdatedFile(entry FileEntry, projectRoot string, assetsFS fs.FS, backupDir string) error {
	content, err := fs.ReadFile(assetsFS, entry.Source)
	if err != nil {
		return fmt.Errorf("reading asset %s: %w", entry.Source, err)
	}

	targetPath := filepath.Join(projectRoot, entry.Target)

	if backupDir != "" {
		if err := backupFile(targetPath, backupDir); err != nil {
			return err
		}
	}

	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

func applyAgentsBlockEntry(entry FileEntry, assetsFS fs.FS, backupDir string) error {
	content, err := fs.ReadFile(assetsFS, entry.Source)
	if err != nil {
		return fmt.Errorf("reading asset %s: %w", entry.Source, err)
	}

	if backupDir != "" {
		if _, statErr := os.Stat(entry.Target); statErr == nil {
			if err := backupFile(entry.Target, backupDir); err != nil {
				return err
			}
		}
	}

	if err := ApplyAgentsBlock(entry.Target, content); err != nil {
		return fmt.Errorf("applying agents block: %w", err)
	}

	return nil
}

func backupFile(targetPath, backupDir string) error {
	data, err := os.ReadFile(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading for backup: %w", err)
	}

	backupPath := filepath.Join(backupDir, filepath.Base(targetPath))
	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return fmt.Errorf("creating backup dir: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("writing backup: %w", err)
	}

	return nil
}
