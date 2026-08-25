package config

import (
	"fmt"
	"strings"
)

type ConfigService struct {
	projectRoot string
	fileStore   *fileConfigStore
	sideEffects *sideEffectRegistry
}

func New(startDir string) (*ConfigService, error) {
	fileStore, err := newFileConfigStore(startDir)
	if err != nil {
		return nil, fmt.Errorf("config service: %w", err)
	}

	return &ConfigService{
		projectRoot: fileStore.ProjectRoot(),
		fileStore:   fileStore,
		sideEffects: newSideEffectRegistry(),
	}, nil
}

func (s *ConfigService) ProjectRoot() string         { return s.projectRoot }
func (s *ConfigService) FileStore() *fileConfigStore { return s.fileStore }

func (s *ConfigService) Get(key string) (string, error) {
	switch {
	case strings.HasPrefix(key, "project."):
		return s.getProjectConfig(key)
	case key == "version_check":
		return fmt.Sprintf("%t", s.fileStore.Config().VersionCheck), nil
	case key == "docs_dir" || key == "docsDir":
		if s.fileStore.Config().DocsDir != "" {
			return s.fileStore.Config().DocsDir, nil
		}
		return "docs", nil
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
	case key == "version_check":
		if err := s.setVersionCheck(value); err != nil {
			return err
		}
	case key == "docs_dir" || key == "docsDir":
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("docs_dir must not be empty")
		}
		s.fileStore.Config().DocsDir = value
		if err := s.fileStore.Save(); err != nil {
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
	docsDir := s.fileStore.Config().DocsDir
	if docsDir == "" {
		docsDir = "docs"
	}
	result["docs_dir"] = docsDir
	for _, p := range s.fileStore.Config().Projects {
		result[fmt.Sprintf("project.%s.wikiRoot", p.ID)] = p.WikiRoot
		result[fmt.Sprintf("project.%s.srcDirs", p.ID)] = fmt.Sprintf("%v", p.SrcDirs)
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

func (s *ConfigService) CacheDir() string {
	return s.fileStore.Config().CacheDir(s.projectRoot)
}

func (s *ConfigService) Close() error {
	return nil
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
