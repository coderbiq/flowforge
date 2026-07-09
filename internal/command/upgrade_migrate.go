package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"flowforge/internal/core"
)

type migration struct {
	name       string
	minVersion string
	run        func(store *core.CardStore, wikiRoot string) error
}

var allMigrations = []migration{
	{
		name:       "v3-wiki-flatten",
		minVersion: "3.0.2",
		run:        migrateV3WikiFlatten,
	},
}

func runPendingMigrations(oldVer, newVer string, wikiRoot string) error {
	store := core.NewCardStore(wikiRoot)
	oldV := strings.TrimPrefix(oldVer, "v")

	var executed int
	for _, m := range allMigrations {
		if compareVersion(oldV, m.minVersion) < 0 {
			fmt.Printf("Running migration: %s (old=%s < required=%s)\n", m.name, oldVer, m.minVersion)
			if err := m.run(store, wikiRoot); err != nil {
				return fmt.Errorf("migration %s failed: %w", m.name, err)
			}
			executed++
		}
	}
	if executed == 0 {
		fmt.Println("No migrations needed.")
	}
	return nil
}

func compareVersion(a, b string) int {
	ap := parseParts(a)
	bp := parseParts(b)
	for i := 0; i < 3 && i < len(ap) && i < len(bp); i++ {
		if ap[i] < bp[i] {
			return -1
		}
		if ap[i] > bp[i] {
			return 1
		}
	}
	return 0
}

func parseParts(v string) []int {
	v = strings.SplitN(v, "-", 2)[0]
	parts := strings.Split(v, ".")
	nums := make([]int, len(parts))
	for i, p := range parts {
		fmt.Sscanf(p, "%d", &nums[i])
	}
	return nums
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
