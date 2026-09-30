package config

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Version != "5.0.0" {
		t.Errorf("expected version 5.0.0, got %s", cfg.Version)
	}

	if cfg.DocsDir != "ff-wiki" {
		t.Errorf("expected default docs_dir ff-wiki, got %s", cfg.DocsDir)
	}
	if len(cfg.Projects) != 0 {
		t.Fatalf("expected no default projects, got %d", len(cfg.Projects))
	}
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	configContent := `version: "2.0.0"
projects:
  - id: "default"
    wikiRoot: "docs"
    srcDirs:
      - "src"
      - "app"
`
	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if len(cfg.Projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(cfg.Projects))
	}

	if cfg.Projects[0].ID != "default" {
		t.Errorf("expected project id default, got %s", cfg.Projects[0].ID)
	}

	if len(cfg.Projects[0].SrcDirs) != 2 {
		t.Fatalf("expected 2 source dirs, got %d", len(cfg.Projects[0].SrcDirs))
	}
}

func TestLoadConfigMissing(t *testing.T) {
	tmpDir := t.TempDir()

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("expected no error for missing config, got: %v", err)
	}

	if cfg.Version != "5.0.0" {
		t.Errorf("expected default version, got %s", cfg.Version)
	}

	if len(cfg.Projects) != 0 {
		t.Fatalf("expected no default projects, got %d", len(cfg.Projects))
	}
}

func TestFindProjectRoot(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte("version: 5.0.0"), 0644); err != nil {
		t.Fatal(err)
	}

	subDir := filepath.Join(tmpDir, "sub", "dir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	root, err := FindProjectRoot(subDir)
	if err != nil {
		t.Fatalf("failed to find project root: %v", err)
	}

	if root != tmpDir {
		t.Errorf("expected root %s, got %s", tmpDir, root)
	}
}

func TestFindProjectRootNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := FindProjectRoot(tmpDir)
	if err == nil {
		t.Error("expected error for missing project root")
	}
}

func TestProjectByID(t *testing.T) {
	cfg := &Config{
		Projects: []ProjectConfig{
			{ID: "frontend", SrcDirs: []string{"ui"}},
			{ID: "backend", SrcDirs: []string{"server"}},
		},
	}

	project, ok := cfg.ProjectByID("backend")
	if !ok {
		t.Fatalf("expected backend project to be found")
	}
	if len(project.SrcDirs) != 1 || project.SrcDirs[0] != "server" {
		t.Fatalf("expected backend srcDirs [server], got %v", project.SrcDirs)
	}

	if _, ok := cfg.ProjectByID("missing"); ok {
		t.Fatalf("expected missing project to not be found")
	}
}

func TestResolveProposalsDir(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	configContent := `version: "5.0.0"
docs_dir: "my-docs"
`
	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	subDir := filepath.Join(tmpDir, "some", "nested", "path")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	resolved, err := ResolveProposalsDir(subDir)
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(tmpDir, "my-docs", "proposals")
	if resolved != expected {
		t.Fatalf("expected %s, got %s", expected, resolved)
	}
}

func TestDocsRootSupportsRelativeAndAbsolutePaths(t *testing.T) {
	projectRoot := t.TempDir()
	relative := Config{DocsDir: "wiki"}
	if got, want := relative.DocsRoot(projectRoot), filepath.Join(projectRoot, "wiki"); got != want {
		t.Fatalf("relative docs root = %s, want %s", got, want)
	}
	absoluteRoot := filepath.Join(t.TempDir(), "external-wiki")
	absolute := Config{DocsDir: absoluteRoot}
	if got := absolute.DocsRoot(projectRoot); got != absoluteRoot {
		t.Fatalf("absolute docs root = %s, want %s", got, absoluteRoot)
	}
}

func TestStandardsGuideDefault(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Standards.Guide != "agents/standards.md" {
		t.Errorf("expected default standards.guide agents/standards.md, got %s", cfg.Standards.Guide)
	}
}

func TestStandardsGuideLoadDefault(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte("version: 5.0.0"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Standards.Guide != "agents/standards.md" {
		t.Errorf("expected standards.guide agents/standards.md, got %s", cfg.Standards.Guide)
	}
}

func TestStandardsGuideSaveLoad(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		Version:      "5.0.0",
		VersionCheck: true,
		DocsDir:      "docs",
		Standards:    StandardsConfig{Guide: "custom/standards.md"},
	}

	if err := cfg.Save(tmpDir); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.Standards.Guide != "custom/standards.md" {
		t.Errorf("expected standards.guide custom/standards.md, got %s", loaded.Standards.Guide)
	}
}

func TestStandardsGuideServiceGet(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte("version: 5.0.0"), 0644); err != nil {
		t.Fatal(err)
	}

	svc, err := New(tmpDir)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	got, err := svc.Get("standards.guide")
	if err != nil {
		t.Fatalf("failed to get standards.guide: %v", err)
	}

	if got != "agents/standards.md" {
		t.Errorf("expected agents/standards.md, got %s", got)
	}
}

func TestStandardsGuideServiceSet(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte("version: 5.0.0"), 0644); err != nil {
		t.Fatal(err)
	}

	svc, err := New(tmpDir)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	if err := svc.Set("standards.guide", "custom/standards.md"); err != nil {
		t.Fatalf("failed to set standards.guide: %v", err)
	}

	got, err := svc.Get("standards.guide")
	if err != nil {
		t.Fatalf("failed to get standards.guide: %v", err)
	}

	if got != "custom/standards.md" {
		t.Errorf("expected custom/standards.md, got %s", got)
	}
}

func TestStandardsGuideServiceSetEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte("version: 5.0.0"), 0644); err != nil {
		t.Fatal(err)
	}

	svc, err := New(tmpDir)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	if err := svc.Set("standards.guide", ""); err == nil {
		t.Error("expected error for empty standards.guide, got nil")
	}
}

func TestStandardsGuideServiceList(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte("version: 5.0.0"), 0644); err != nil {
		t.Fatal(err)
	}

	svc, err := New(tmpDir)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	values, err := svc.List()
	if err != nil {
		t.Fatalf("failed to list config: %v", err)
	}

	got, ok := values["standards.guide"]
	if !ok {
		t.Fatal("expected standards.guide in config list")
	}

	if got != "agents/standards.md" {
		t.Errorf("expected agents/standards.md, got %s", got)
	}
}

func TestLoadWarnsOnDeprecatedWikiKeys(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	configContent := `version: "5.0.0"
docs_dir: "my-docs"
wiki:
  root: "legacy-wiki"
projects:
  - id: "default"
    wikiRoot: "project-wiki"
    srcDirs:
      - "src"
`
	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	orig := warnOut
	warnOut = &buf
	t.Cleanup(func() { warnOut = orig })

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("config with legacy wiki keys must load without error, got: %v", err)
	}

	want := "warning: config key \"wiki.root\" is deprecated and ignored; wiki root is decided by \"docs_dir\" only\n" +
		"warning: config key \"projects[default].wikiRoot\" is deprecated and ignored; use \"docs_dir\"\n"
	if buf.String() != want {
		t.Fatalf("unexpected warnings:\n%s\nwant:\n%s", buf.String(), want)
	}

	if got, want := cfg.DocsRoot(tmpDir), filepath.Join(tmpDir, "my-docs"); got != want {
		t.Fatalf("docs root = %s, want %s", got, want)
	}
}

func TestLoadIgnoresLegacyWikiBlock(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	configContent := `version: "5.0.0"
docs_dir: "my-docs"
wiki:
  root: "legacy-wiki"
projects:
  - id: "default"
    wikiRoot: "ignored-wiki"
    srcDirs:
      - "src"
`
	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	orig := warnOut
	warnOut = &buf
	t.Cleanup(func() { warnOut = orig })

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("legacy wiki config must load without error, got: %v", err)
	}

	if got, want := cfg.DocsRoot(tmpDir), filepath.Join(tmpDir, "my-docs"); got != want {
		t.Fatalf("docs root = %s, want %s", got, want)
	}

	if len(cfg.Projects) != 1 || cfg.Projects[0].ID != "default" {
		t.Fatalf("expected project default to survive load, got %+v", cfg.Projects)
	}
}

func TestConfigSetProjectWikiRootRejected(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	configContent := `version: "5.0.0"
projects:
  - id: "default"
    srcDirs:
      - "src"
`
	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	svc, err := New(tmpDir)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	err = svc.Set("project.default.wikiRoot", "some-wiki")
	if err == nil {
		t.Fatal("expected error setting project wikiRoot, got nil")
	}
	if !strings.Contains(err.Error(), "unknown project config field: wikiRoot") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadSaveModelObjects(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ConfigDirName), 0755); err != nil {
		t.Fatal(err)
	}
	data := "agents:\n  models_by_name:\n    flowforge-investigator: {model: provider/model, reasoning_effort: high}\n"
	if err := os.WriteFile(ConfigPath(root), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg.Save(root); err != nil {
		t.Fatal(err)
	}
	again, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Agents, again.Agents) {
		t.Fatalf("roundtrip: %#v != %#v", cfg.Agents, again.Agents)
	}
}

func TestLoadModelValueStrictPaths(t *testing.T) {
	for _, raw := range []string{"null", "42", "true", "[]", "{}", "{model: null}", "{model: 42}", "{model: true}", "{model: ''}", "{reasoning_effort: null}", "{reasoning_effort: 42}", "{reasoning_effort: false}", "{reasoning_effort: ''}", "{model: provider/model, unknown: high}"} {
		for _, layer := range []string{"models_by_name", "model_sets.quick.models_by_name"} {
			t.Run(layer+raw, func(t *testing.T) {
				root := t.TempDir()
				if err := os.MkdirAll(filepath.Join(root, ConfigDirName), 0755); err != nil {
					t.Fatal(err)
				}
				prefix := "agents:\n  models_by_name:\n    flowforge-investigator: "
				if layer != "models_by_name" {
					prefix = "agents:\n  model_sets:\n    quick:\n      models_by_name:\n        flowforge-investigator: "
				}
				if err := os.WriteFile(ConfigPath(root), []byte(prefix+raw+"\n"), 0644); err != nil {
					t.Fatal(err)
				}
				_, err := Load(root)
				path := "agents." + layer + ".flowforge-investigator"
				if err == nil || !strings.Contains(err.Error(), path) {
					t.Fatalf("want strict error at %s, got %v", path, err)
				}
			})
		}
	}
}

func TestLoadSaveIndependentModelFields(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ConfigDirName), 0755); err != nil {
		t.Fatal(err)
	}
	data := `agents:
  models:
    tool-capable: provider/base
  models_by_name:
    flowforge-investigator: {reasoning_effort: high}
    flowforge-reviewer: {model: provider/reviewer}
  models_by_host:
    codex:
      flowforge-investigator: {model: gpt-codex, reasoning_effort: inherit}
  model_sets:
    quick:
      models_by_name:
        flowforge-reviewer: {reasoning_effort: low}
`
	if err := os.WriteFile(ConfigPath(root), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := cfg.Agents.ModelOverrides["flowforge-investigator"]; got.Model != "" || got.ReasoningEffort != "high" {
		t.Fatalf("effort only: %#v", got)
	}
	if err := cfg.Save(root); err != nil {
		t.Fatal(err)
	}
	again, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Agents, again.Agents) {
		t.Fatalf("independent field roundtrip changed config")
	}
	var decoded Config
	body, err := os.ReadFile(ConfigPath(root))
	if err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal(body, &decoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded.Agents.Models, cfg.Agents.Models) || !reflect.DeepEqual(decoded.Agents.ModelOverrides, cfg.Agents.ModelOverrides) || !reflect.DeepEqual(decoded.Agents.ModelHostOverrides, cfg.Agents.ModelHostOverrides) || !reflect.DeepEqual(decoded.Agents.ModelSets, cfg.Agents.ModelSets) {
		t.Fatal("direct YAML decoder diverges")
	}
}

func TestDirectYAMLRejectsNullModelFields(t *testing.T) {
	for _, raw := range []string{"null", "{model: null}", "{reasoning_effort: false}", "{unknown: high}"} {
		var cfg Config
		err := yaml.Unmarshal([]byte("agents:\n  models_by_name:\n    flowforge-investigator: "+raw+"\n"), &cfg)
		if err == nil {
			t.Fatalf("direct YAML accepted %s", raw)
		}
	}
}
