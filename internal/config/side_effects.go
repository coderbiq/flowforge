package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

type sideEffectFunc func(svc *ConfigService, oldValue, newValue string) error

type sideEffectRegistry struct {
	effects []sideEffectEntry
}

type sideEffectEntry struct {
	pattern string
	fn      sideEffectFunc
}

func newSideEffectRegistry() *sideEffectRegistry {
	r := &sideEffectRegistry{}
	r.register("project.*.wikiRoot", func(svc *ConfigService, oldValue, newValue string) error {
		if oldValue == newValue {
			return nil
		}
		return svc.rebuildIndexForWikiRoot(newValue)
	})
	return r
}

func (r *sideEffectRegistry) register(pattern string, fn sideEffectFunc) {
	r.effects = append(r.effects, sideEffectEntry{pattern: pattern, fn: fn})
}

func (r *sideEffectRegistry) trigger(svc *ConfigService, key, oldValue, newValue string) error {
	for _, e := range r.effects {
		if matchPattern(e.pattern, key) {
			if err := e.fn(svc, oldValue, newValue); err != nil {
				return fmt.Errorf("side effect %q failed: %w", e.pattern, err)
			}
		}
	}
	return nil
}

func (s *ConfigService) rebuildIndexForWikiRoot(wikiRoot string) error {
	wikiPath := wikiRoot
	if !filepath.IsAbs(wikiPath) {
		wikiPath = filepath.Join(s.projectRoot, wikiPath)
	}
	return rebuildIndex(s.projectRoot, wikiPath)
}

func matchPattern(pattern, key string) bool {
	pp := strings.Split(pattern, ".")
	kp := strings.Split(key, ".")
	if len(pp) != len(kp) {
		return false
	}
	for i := range pp {
		if pp[i] == "*" {
			continue
		}
		if pp[i] != kp[i] {
			return false
		}
	}
	return true
}

func rebuildIndex(projectRoot, wikiRoot string) error {
	dbPath := filepath.Join(projectRoot, ConfigDirName, "cache", "flowforge.sqlite")
	store, err := state.Open(dbPath)
	if err != nil {
		return fmt.Errorf("opening state for index rebuild: %w", err)
	}
	defer store.Close()

	syncSvc := state.NewCardSyncService(store.DB())
	cardStore := core.NewCardStore(wikiRoot)
	dirs := []string{
		cardStore.ActiveDir(), cardStore.IntakeDir(), cardStore.CompletedDir(),
		cardStore.LibraryDir(), cardStore.ProposalCardDir(),
	}
	_, _, err = syncSvc.RebuildAll(cardStore.ListCardsFromFiles, dirs)
	return err
}