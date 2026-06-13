package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ConfigDirName  = ".flowforge"
	ConfigFileName = "config.yaml"
)

type Config struct {
	Version  string          `yaml:"version" mapstructure:"version"`
	Projects []ProjectConfig `yaml:"projects" mapstructure:"projects"`
	Wiki     WikiConfig      `yaml:"wiki" mapstructure:"wiki"`
}

type ProjectConfig struct {
	ID       string   `yaml:"id" mapstructure:"id"`
	WikiRoot string   `yaml:"wikiRoot" mapstructure:"wikiRoot"`
	SrcDirs  []string `yaml:"srcDirs" mapstructure:"srcDirs"`
}

type WikiConfig struct {
	Root string `yaml:"root" mapstructure:"root"`
}

var defaultConfig = Config{
	Version: "2.0.0",
	Wiki: WikiConfig{
		Root: "ff-wiki",
	},
}

func DefaultConfig() Config {
	return defaultConfig
}

func ConfigPath(projectRoot string) string {
	return filepath.Join(projectRoot, ConfigDirName, ConfigFileName)
}

func (c *Config) Save(projectRoot string) error {
	if c == nil {
		return fmt.Errorf("config is required")
	}

	type fileConfig struct {
		Version  string          `yaml:"version"`
		Projects []ProjectConfig `yaml:"projects"`
	}

	payload := fileConfig{
		Version:  c.Version,
		Projects: c.Projects,
	}

	data, err := yaml.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	configDir := filepath.Join(projectRoot, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	configPath := ConfigPath(projectRoot)
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

func FindProjectRoot(startDir string) (string, error) {
	dir := startDir
	for {
		configPath := filepath.Join(dir, ConfigDirName, ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no .flowforge/config.yaml found in %s or parents", startDir)
		}
		dir = parent
	}
}

func Load(projectRoot string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(projectRoot, ConfigDirName))
	v.SetConfigName("config")

	v.SetDefault("version", defaultConfig.Version)
	v.SetDefault("projects", defaultConfig.Projects)
	v.SetDefault("wiki.root", defaultConfig.Wiki.Root)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cfg := defaultConfig
			return &cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) WikiRoot(projectRoot string) string {
	project := c.primaryProject()
	return c.projectWikiRoot(projectRoot, project)
}

func (c *Config) ProjectByID(id string) (ProjectConfig, bool) {
	for _, project := range c.Projects {
		if project.ID == id {
			return project, true
		}
	}

	return ProjectConfig{}, false
}

func (c *Config) WikiRootForProject(projectRoot string, projectID string) (string, error) {
	project, ok := c.ProjectByID(projectID)
	if !ok {
		return "", fmt.Errorf("project %q is not registered", projectID)
	}

	return c.projectWikiRoot(projectRoot, project), nil
}

func (c *Config) projectWikiRoot(projectRoot string, project ProjectConfig) string {
	if project.WikiRoot != "" {
		if filepath.IsAbs(project.WikiRoot) {
			return project.WikiRoot
		}
		return filepath.Join(projectRoot, project.WikiRoot)
	}

	if filepath.IsAbs(c.Wiki.Root) {
		return c.Wiki.Root
	}

	if c.Wiki.Root != "" {
		return filepath.Join(projectRoot, c.Wiki.Root)
	}

	return filepath.Join(projectRoot, "ff-wiki")
}

func (c *Config) primaryProject() ProjectConfig {
	if len(c.Projects) > 0 {
		return c.Projects[0]
	}

	return ProjectConfig{
		WikiRoot: c.Wiki.Root,
		SrcDirs:  nil,
	}
}

func (c *Config) ConfigDir(projectRoot string) string {
	return filepath.Join(projectRoot, ConfigDirName)
}

func (c *Config) CacheDir(projectRoot string) string {
	return filepath.Join(c.ConfigDir(projectRoot), "cache")
}
