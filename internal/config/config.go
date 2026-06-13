package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
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
	Projects: []ProjectConfig{
		{
			ID:       "default",
			WikiRoot: "ff-wiki",
			SrcDirs:  []string{},
		},
	},
	Wiki: WikiConfig{
		Root: "ff-wiki",
	},
}

func DefaultConfig() Config {
	return defaultConfig
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
