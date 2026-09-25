package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ConfigDirName  = ".flowforge"
	ConfigFileName = "config.yaml"

	DefaultDocsDir        = "ff-wiki"
	DefaultStandardsGuide = "agents/standards.md"
)

type Config struct {
	Version          string                  `yaml:"version" mapstructure:"version"`
	VersionCheck     bool                    `yaml:"version_check" mapstructure:"version_check"`
	DocsDir          string                  `yaml:"docs_dir,omitempty" mapstructure:"docs_dir"`
	Projects         []ProjectConfig         `yaml:"projects,omitempty" mapstructure:"projects"`
	KnowledgeSources []KnowledgeSourceConfig `yaml:"knowledge_sources,omitempty" mapstructure:"knowledge_sources"`
	Agents           AgentsConfig            `yaml:"agents,omitempty" mapstructure:"agents"`
	Standards        StandardsConfig         `yaml:"standards,omitempty" mapstructure:"standards"`
	Evidence         EvidenceConfig          `yaml:"evidence,omitempty" mapstructure:"evidence"`
}

type StandardsConfig struct {
	Guide string `yaml:"guide,omitempty" mapstructure:"guide"`
}

type EvidenceConfig struct {
	ExemptProposals []string `yaml:"exempt_proposals,omitempty" mapstructure:"exempt_proposals"`
}

type AgentsConfig struct {
	Disabled           []string                     `yaml:"disabled,omitempty" mapstructure:"disabled"`
	Hosts              []string                     `yaml:"hosts,omitempty" mapstructure:"hosts"`
	MaxSteps           int                          `yaml:"max_steps,omitempty" mapstructure:"max_steps"`
	Models             map[string]string            `yaml:"models,omitempty" mapstructure:"models"`
	ModelOverrides     map[string]string            `yaml:"models_by_name,omitempty" mapstructure:"models_by_name"`
	ModelHostOverrides map[string]map[string]string `yaml:"models_by_host,omitempty" mapstructure:"models_by_host"`
	TestFileGlobs      []string                     `yaml:"test_file_globs,omitempty" mapstructure:"test_file_globs"`
	DisableTestGuard   bool                         `yaml:"disable_test_guard,omitempty" mapstructure:"disable_test_guard"`
}

type ProjectConfig struct {
	ID      string   `yaml:"id" mapstructure:"id"`
	SrcDirs []string `yaml:"srcDirs" mapstructure:"srcDirs"`
}

type KnowledgeSourceConfig struct {
	Name        string `yaml:"name" mapstructure:"name"`
	Path        string `yaml:"path" mapstructure:"path"`
	Type        string `yaml:"type" mapstructure:"type"`
	Category    string `yaml:"category" mapstructure:"category"`
	Trust       string `yaml:"trust" mapstructure:"trust"`
	Description string `yaml:"description" mapstructure:"description"`
}

var defaultConfig = Config{
	Version:      "5.0.0",
	VersionCheck: true,
	DocsDir:      DefaultDocsDir,
	Standards: StandardsConfig{
		Guide: DefaultStandardsGuide,
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
		Version          string                  `yaml:"version"`
		VersionCheck     bool                    `yaml:"version_check"`
		DocsDir          string                  `yaml:"docs_dir,omitempty"`
		Projects         []ProjectConfig         `yaml:"projects,omitempty"`
		KnowledgeSources []KnowledgeSourceConfig `yaml:"knowledge_sources,omitempty"`
		Agents           AgentsConfig            `yaml:"agents,omitempty"`
		Standards        StandardsConfig         `yaml:"standards,omitempty"`
		Evidence         EvidenceConfig          `yaml:"evidence,omitempty"`
	}

	payload := fileConfig{
		Version:          c.Version,
		VersionCheck:     c.VersionCheck,
		DocsDir:          c.DocsDir,
		Projects:         c.Projects,
		KnowledgeSources: c.KnowledgeSources,
		Agents:           c.Agents,
		Standards:        c.Standards,
		Evidence:         c.Evidence,
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
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("resolving start directory: %w", err)
	}
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

// warnOut is the destination for deprecated-config warnings emitted by Load.
// It is a package-level seam so tests can replace and capture the output.
var warnOut io.Writer = os.Stderr

func Load(projectRoot string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(projectRoot, ConfigDirName))
	v.SetConfigName("config")

	v.SetDefault("version", defaultConfig.Version)
	v.SetDefault("version_check", defaultConfig.VersionCheck)
	v.SetDefault("docs_dir", defaultConfig.DocsDir)
	v.SetDefault("projects", defaultConfig.Projects)
	v.SetDefault("knowledge_sources", []KnowledgeSourceConfig{})
	v.SetDefault("agents.disabled", []string{})
	v.SetDefault("standards.guide", defaultConfig.Standards.Guide)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cfg := defaultConfig
			return &cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	warnDeprecatedWikiKeys(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}
	return &cfg, nil
}

// warnDeprecatedWikiKeys reports legacy wiki-track keys still present in the
// raw viper config. The wiki track is removed; docs_dir is the single source
// of truth for the wiki/docs root, so these keys are ignored after warning.
func warnDeprecatedWikiKeys(v *viper.Viper) {
	if v.IsSet("wiki") {
		fmt.Fprintf(warnOut, "warning: config key %q is deprecated and ignored; wiki root is decided by %q only\n", "wiki.root", "docs_dir")
	}
	rawProjects, ok := v.Get("projects").([]any)
	if !ok {
		return
	}
	// Viper lowercases all config keys (including keys inside slices), so the
	// raw maps expose "wikiroot"/"id" in their normalized form.
	for _, raw := range rawProjects {
		project, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		wikiRoot, _ := project["wikiroot"].(string)
		if wikiRoot == "" {
			continue
		}
		id, _ := project["id"].(string)
		fmt.Fprintf(warnOut, "warning: config key %q is deprecated and ignored; use %q\n", fmt.Sprintf("projects[%s].wikiRoot", id), "docs_dir")
	}
}

func (c *Config) DocsRoot(projectRoot string) string {
	if c.DocsDir != "" {
		if filepath.IsAbs(c.DocsDir) {
			return c.DocsDir
		}
		return filepath.Join(projectRoot, c.DocsDir)
	}
	return filepath.Join(projectRoot, DefaultDocsDir)
}

func (c *Config) ProposalsDir(projectRoot string) string {
	return filepath.Join(c.DocsRoot(projectRoot), "proposals")
}

func ResolveProposalsDir(startDir string) (string, error) {
	projectRoot, err := FindProjectRoot(startDir)
	if err != nil {
		return filepath.Join(startDir, DefaultDocsDir, "proposals"), nil
	}

	cfg, err := Load(projectRoot)
	if err != nil {
		return "", fmt.Errorf("loading project configuration: %w", err)
	}

	return cfg.ProposalsDir(projectRoot), nil
}

func (c *Config) ProjectByID(id string) (ProjectConfig, bool) {
	for _, project := range c.Projects {
		if project.ID == id {
			return project, true
		}
	}

	return ProjectConfig{}, false
}

func (c *Config) ConfigDir(projectRoot string) string {
	return filepath.Join(projectRoot, ConfigDirName)
}

func (c *Config) CacheDir(projectRoot string) string {
	return filepath.Join(c.ConfigDir(projectRoot), "cache")
}
