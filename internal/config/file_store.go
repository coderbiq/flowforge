package config

import "fmt"

type fileConfigStore struct {
	projectRoot string
	cfg         *Config
}

func newFileConfigStore(startDir string) (*fileConfigStore, error) {
	projectRoot, err := FindProjectRoot(startDir)
	if err != nil {
		return nil, fmt.Errorf("finding project root: %w", err)
	}
	cfg, err := Load(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	return &fileConfigStore{projectRoot: projectRoot, cfg: cfg}, nil
}

func (s *fileConfigStore) ProjectRoot() string { return s.projectRoot }

func (s *fileConfigStore) Config() *Config { return s.cfg }

func (s *fileConfigStore) Save() error {
	if s.cfg == nil {
		return fmt.Errorf("config is nil")
	}
	return s.cfg.Save(s.projectRoot)
}

func (s *fileConfigStore) Reload() error {
	cfg, err := Load(s.projectRoot)
	if err != nil {
		return fmt.Errorf("reloading config: %w", err)
	}
	s.cfg = cfg
	return nil
}