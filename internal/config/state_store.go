package config

import (
	"fmt"

	"flowforge/internal/state"
)

type runtimeStateStore struct {
	store *state.Store
}

func newRuntimeStateStore(dbPath string) (*runtimeStateStore, error) {
	store, err := state.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening state store: %w", err)
	}
	if err := store.EnsureSchema(); err != nil {
		store.Close()
		return nil, fmt.Errorf("ensuring schema: %w", err)
	}
	return &runtimeStateStore{store: store}, nil
}

func (s *runtimeStateStore) CurrentProjectID() (string, bool, error) {
	return s.store.CurrentProjectID()
}

func (s *runtimeStateStore) SetCurrentProjectID(id string) error {
	return s.store.SetCurrentProjectID(id)
}

func (s *runtimeStateStore) ClearCurrentProjectID() error {
	return s.store.ClearCurrentProjectID()
}

func (s *runtimeStateStore) CurrentProposalID(projectID string) (string, bool, error) {
	return s.store.CurrentProposalID(projectID)
}

func (s *runtimeStateStore) SetCurrentProposalID(projectID, proposalID string) error {
	return s.store.SetCurrentProposalID(projectID, proposalID)
}

func (s *runtimeStateStore) ClearCurrentProposalID(projectID string) error {
	return s.store.ClearCurrentProposalID(projectID)
}

func (s *runtimeStateStore) DB() *state.Store {
	return s.store
}

func (s *runtimeStateStore) Close() error {
	return s.store.Close()
}