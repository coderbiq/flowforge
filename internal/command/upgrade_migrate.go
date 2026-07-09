package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"flowforge/internal/core"
)

type migration struct {
	name string
	run  func(store *core.CardStore, wikiRoot string) error
}

func runPendingMigrations(oldVer, newVer string, wikiRoot string) error {
	if !needsV3Migration(oldVer, newVer) {
		fmt.Println("No migrations needed.")
		return nil
	}

	store := core.NewCardStore(wikiRoot)
	for _, m := range v3Migrations() {
		fmt.Printf("Running migration: %s\n", m.name)
		if err := m.run(store, wikiRoot); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.name, err)
		}
	}
	return nil
}

func needsV3Migration(oldVer, newVer string) bool {
	oldV := strings.TrimPrefix(oldVer, "v")
	newV := strings.TrimPrefix(newVer, "v")

	oldMajor, _ := parseMajor(oldV)
	newMajor, _ := parseMajor(newV)

	return oldMajor < 3 && newMajor >= 3
}

func parseMajor(v string) (int, error) {
	parts := strings.SplitN(v, ".", 2)
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid version: %s", v)
	}
	var major int
	_, err := fmt.Sscanf(parts[0], "%d", &major)
	return major, err
}

func v3Migrations() []migration {
	return []migration{
		{
			name: "v3-wiki-flatten",
			run:  migrateV3WikiFlatten,
		},
	}
}

func migrateV3WikiFlatten(store *core.CardStore, wikiRoot string) error {
	workspaceDir := filepath.Join(wikiRoot, "01-workspace")

	oldDirs := []string{
		filepath.Join(workspaceDir, "01-active"),
		filepath.Join(workspaceDir, "02-intake"),
		filepath.Join(workspaceDir, "03-completed"),
	}

	hasOldStructure := false
	for _, dir := range oldDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			hasOldStructure = true
			break
		}
	}

	if !hasOldStructure {
		fmt.Println("  Already v3 flat structure — skipping wiki migration.")
		return nil
	}

	var moved int
	var setCompleted int

	activeDir := filepath.Join(workspaceDir, "01-active")
	if entries, err := os.ReadDir(activeDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			src := filepath.Join(activeDir, entry.Name())
			dst := filepath.Join(workspaceDir, entry.Name())
			if err := os.Rename(src, dst); err != nil {
				return fmt.Errorf("moving proposal %s: %w", entry.Name(), err)
			}
			moved++
		}
	}

	completedDir := filepath.Join(workspaceDir, "03-completed")
	if entries, err := os.ReadDir(completedDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			src := filepath.Join(completedDir, entry.Name())
			dst := filepath.Join(workspaceDir, entry.Name())
			if _, err := os.Stat(dst); err == nil {
				fmt.Printf("  Skipping %s (already exists at target)\n", entry.Name())
				continue
			}
			if err := os.Rename(src, dst); err != nil {
				return fmt.Errorf("moving completed proposal %s: %w", entry.Name(), err)
			}
			moved++

			propID := "PROP-" + entry.Name()
			card, err := store.ReadCard(propID)
			if err == nil {
				card.Status = core.CardStatusCompleted
				if err := store.UpdateCard(card); err != nil {
					fmt.Printf("  Warning: could not update %s status: %v\n", propID, err)
				} else {
					setCompleted++
				}
			}
		}
	}

	for _, dir := range oldDirs {
		if isEmptyDir(dir) {
			os.Remove(dir)
		}
	}

	fmt.Printf("  Wiki migration complete: %d proposals moved, %d marked completed.\n", moved, setCompleted)
	return nil
}

func isEmptyDir(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) == 0
}
