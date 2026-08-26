package command

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

type managedAssetState string

const (
	managedAssetCurrent      managedAssetState = "current"
	managedAssetMissing      managedAssetState = "missing"
	managedAssetDrifted      managedAssetState = "drifted"
	managedAssetProjectOwned managedAssetState = "project-owned"
)

type managedAssetEntry struct {
	State      managedAssetState
	SourcePath string
	TargetPath string
}

type managedAssetComparison struct {
	Entries []managedAssetEntry
}

func (c managedAssetComparison) IsCurrent() bool {
	for _, entry := range c.Entries {
		if entry.State == managedAssetMissing || entry.State == managedAssetDrifted {
			return false
		}
	}
	return true
}

func compareManagedAssets(assetsDir, targetDir, docsRoot string) (managedAssetComparison, error) {
	if docsRoot == "" {
		docsRoot = filepath.Join(targetDir, "docs")
	} else if !filepath.IsAbs(docsRoot) {
		docsRoot = filepath.Join(targetDir, docsRoot)
	}

	targets := []struct {
		source string
		target string
	}{
		{source: filepath.Join(assetsDir, "skills"), target: filepath.Join(targetDir, ".agents", "skills")},
		{source: filepath.Join(assetsDir, "agents"), target: filepath.Join(docsRoot, "agents")},
	}

	comparison := managedAssetComparison{}
	for _, target := range targets {
		entries, err := compareManagedAssetTree(target.source, target.target)
		if err != nil {
			return managedAssetComparison{}, err
		}
		comparison.Entries = append(comparison.Entries, entries...)
	}
	sort.Slice(comparison.Entries, func(i, j int) bool {
		return comparison.Entries[i].TargetPath < comparison.Entries[j].TargetPath
	})
	return comparison, nil
}

func compareManagedAssetTree(sourceDir, targetDir string) ([]managedAssetEntry, error) {
	managed := map[string]string{}
	if err := filepath.WalkDir(sourceDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || entry.Name() == ".gitkeep" {
			return nil
		}
		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		managed[filepath.Clean(rel)] = path
		return nil
	}); err != nil {
		return nil, fmt.Errorf("discovering embedded asset tree %s: %w", sourceDir, err)
	}

	entries := make([]managedAssetEntry, 0, len(managed))
	for rel, sourcePath := range managed {
		targetPath := filepath.Join(targetDir, rel)
		state, err := compareManagedAsset(sourcePath, targetPath)
		if err != nil {
			return nil, err
		}
		entries = append(entries, managedAssetEntry{State: state, SourcePath: sourcePath, TargetPath: targetPath})
	}
	if err := filepath.WalkDir(targetDir, func(path string, entry fs.DirEntry, err error) error {
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(targetDir, path)
		if err != nil {
			return err
		}
		if _, known := managed[filepath.Clean(rel)]; !known {
			entries = append(entries, managedAssetEntry{State: managedAssetProjectOwned, TargetPath: path})
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("discovering project asset tree %s: %w", targetDir, err)
	}
	return entries, nil
}

func compareManagedAsset(sourcePath, targetPath string) (managedAssetState, error) {
	source, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", fmt.Errorf("reading embedded asset %s: %w", sourcePath, err)
	}
	target, err := os.ReadFile(targetPath)
	if os.IsNotExist(err) {
		return managedAssetMissing, nil
	}
	if err != nil {
		return "", fmt.Errorf("reading project asset %s: %w", targetPath, err)
	}
	if bytes.Equal(source, target) {
		return managedAssetCurrent, nil
	}
	return managedAssetDrifted, nil
}
