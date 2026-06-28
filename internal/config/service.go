package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

type ConfigService struct {
	projectRoot  string
	fileStore    *fileConfigStore
	stateStore   *runtimeStateStore
	sideEffects  *sideEffectRegistry
}

func New(startDir string) (*ConfigService, error) {
	fileStore, err := newFileConfigStore(startDir)
	if err != nil {
		return nil, fmt.Errorf("config service: %w", err)
	}

	dbPath := filepath.Join(fileStore.Config().CacheDir(fileStore.ProjectRoot()), "flowforge.sqlite")
	stateStore, err := newRuntimeStateStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("config service: %w", err)
	}

	return &ConfigService{
		projectRoot: fileStore.ProjectRoot(),
		fileStore:   fileStore,
		stateStore:  stateStore,
		sideEffects: newSideEffectRegistry(),
	}, nil
}

func (s *ConfigService) ProjectRoot() string       { return s.projectRoot }
func (s *ConfigService) FileStore() *fileConfigStore { return s.fileStore }
func (s *ConfigService) StateStore() *runtimeStateStore { return s.stateStore }

func (s *ConfigService) Get(key string) (string, error) {
	switch {
	case strings.HasPrefix(key, "project."):
		return s.getProjectConfig(key)
	case strings.HasPrefix(key, "runtime."):
		return s.getRuntimeState(key)
	case key == "version_check":
		return fmt.Sprintf("%t", s.fileStore.Config().VersionCheck), nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}

func (s *ConfigService) Set(key, value string) error {
	old, _ := s.Get(key)
	switch {
	case strings.HasPrefix(key, "project."):
		if err := s.setProjectConfig(key, value); err != nil {
			return err
		}
	case strings.HasPrefix(key, "runtime."):
		if err := s.setRuntimeState(key, value); err != nil {
			return err
		}
	case key == "version_check":
		if err := s.setVersionCheck(value); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return s.sideEffects.trigger(s, key, old, value)
}

func (s *ConfigService) List() (map[string]string, error) {
	result := make(map[string]string)
	result["version_check"] = fmt.Sprintf("%t", s.fileStore.Config().VersionCheck)
	for _, p := range s.fileStore.Config().Projects {
		result[fmt.Sprintf("project.%s.wikiRoot", p.ID)] = p.WikiRoot
		result[fmt.Sprintf("project.%s.srcDirs", p.ID)] = fmt.Sprintf("%v", p.SrcDirs)
	}
	if id, ok, _ := s.stateStore.CurrentProjectID(); ok {
		result["runtime.currentProjectId"] = id
	}
	if id, ok, _ := s.stateStore.CurrentProjectID(); ok {
		if pid, ok2, _ := s.stateStore.CurrentProposalID(id); ok2 {
			result[fmt.Sprintf("runtime.currentProposalId.%s", id)] = pid
		}
	}
	return result, nil
}

func (s *ConfigService) Projects() []ProjectConfig {
	return s.fileStore.Config().Projects
}

func (s *ConfigService) ProjectByID(id string) (ProjectConfig, error) {
	p, ok := s.fileStore.Config().ProjectByID(id)
	if !ok {
		return ProjectConfig{}, fmt.Errorf("project %q not found", id)
	}
	return p, nil
}

func (s *ConfigService) WikiRoot(projectID string) (string, error) {
	return s.fileStore.Config().WikiRootForProject(s.projectRoot, projectID)
}

func (s *ConfigService) CurrentProjectID() (string, bool, error) {
	return s.stateStore.CurrentProjectID()
}

func (s *ConfigService) SetCurrentProjectID(id string) error {
	return s.stateStore.SetCurrentProjectID(id)
}

func (s *ConfigService) CurrentProposalID(projectID string) (string, bool, error) {
	return s.stateStore.CurrentProposalID(projectID)
}

func (s *ConfigService) SetCurrentProposalID(projectID, proposalID string) error {
	return s.stateStore.SetCurrentProposalID(projectID, proposalID)
}

func (s *ConfigService) CacheDir() string {
	return s.fileStore.Config().CacheDir(s.projectRoot)
}

func (s *ConfigService) Close() error {
	return s.stateStore.Close()
}

func (s *ConfigService) getProjectConfig(key string) (string, error) {
	cfg := s.fileStore.Config()
	parts := strings.SplitN(key, ".", 3)
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid project config key: %s", key)
	}
	projectID := parts[1]
	field := parts[2]

	p, ok := cfg.ProjectByID(projectID)
	if !ok {
		return "", fmt.Errorf("project %q not found", projectID)
	}
	switch field {
	case "wikiRoot":
		return p.WikiRoot, nil
	case "srcDirs":
		return fmt.Sprintf("%v", p.SrcDirs), nil
	default:
		return "", fmt.Errorf("unknown project config field: %s", field)
	}
}

func (s *ConfigService) setProjectConfig(key, value string) error {
	cfg := s.fileStore.Config()
	parts := strings.SplitN(key, ".", 3)
	if len(parts) < 3 {
		return fmt.Errorf("invalid project config key: %s", key)
	}
	projectID := parts[1]
	field := parts[2]

	for i := range cfg.Projects {
		if cfg.Projects[i].ID == projectID {
			switch field {
			case "wikiRoot":
				cfg.Projects[i].WikiRoot = value
			case "srcDirs":
				cfg.Projects[i].SrcDirs = []string{value}
			default:
				return fmt.Errorf("unknown project config field: %s", field)
			}
			return s.fileStore.Save()
		}
	}
	return fmt.Errorf("project %q not found", projectID)
}

func (s *ConfigService) setVersionCheck(value string) error {
	enabled := true
	switch value {
	case "false", "0", "no":
		enabled = false
	case "true", "1", "yes":
		enabled = true
	default:
		return fmt.Errorf("version_check must be true or false, got: %s", value)
	}
	s.fileStore.Config().VersionCheck = enabled
	return s.fileStore.Save()
}

func (s *ConfigService) getRuntimeState(key string) (string, error) {
	key = strings.TrimPrefix(key, "runtime.")
	switch {
	case key == "currentProjectId":
		id, ok, err := s.stateStore.CurrentProjectID()
		if err != nil {
			return "", err
		}
		if !ok {
			return "", fmt.Errorf("no current project set")
		}
		return id, nil
	case strings.HasPrefix(key, "currentProposalId."):
		projectID := strings.TrimPrefix(key, "currentProposalId.")
		id, ok, err := s.stateStore.CurrentProposalID(projectID)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", fmt.Errorf("no current proposal set for project %s", projectID)
		}
		return id, nil
	default:
		return "", fmt.Errorf("unknown runtime state key: %s", key)
	}
}

func (s *ConfigService) setRuntimeState(key, value string) error {
	key = strings.TrimPrefix(key, "runtime.")
	switch {
	case key == "currentProjectId":
		return s.stateStore.SetCurrentProjectID(value)
	case strings.HasPrefix(key, "currentProposalId."):
		projectID := strings.TrimPrefix(key, "currentProposalId.")
		return s.stateStore.SetCurrentProposalID(projectID, value)
	default:
		return fmt.Errorf("unknown runtime state key: %s", key)
	}
}