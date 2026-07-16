package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const ManifestFileName = "manifest.yaml"
const configDirName = ".flowforge"

type FileEntry struct {
	Source   string            `yaml:"source"`
	Target   string            `yaml:"target"`
	SHA256   string            `yaml:"sha256"`
	Type     string            `yaml:"type"`
	Markers  *BlockMarkers     `yaml:"markers,omitempty"`
}

type BlockMarkers struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

type ProjectManifest struct {
	Version    int         `yaml:"version"`
	CLIVersion string      `yaml:"cli_version"`
	Files      []FileEntry `yaml:"files"`
}

type DiffResult struct {
	Added    []FileEntry
	Updated  []FileEntry
	Conflict []FileEntry
	Removed  []FileEntry
}

func LoadProjectManifest(projectRoot string) (*ProjectManifest, error) {
	path := filepath.Join(projectRoot, configDirName, ManifestFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest: %w", err)
	}

	var m ProjectManifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest: %w", err)
	}

	return &m, nil
}

func (m *ProjectManifest) Save(projectRoot string) error {
	path := filepath.Join(projectRoot, configDirName, ManifestFileName)
	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating manifest dir: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	return nil
}

var assetMappings = []struct {
	sourceDir string
	targetDir string
	fileType  string
}{
	{"assets/skills", ".agents/skills", "skill"},
	{"assets/templates", ".flowforge/templates", "template"},
	{"assets/wiki", "", "wiki"},
}

func GenerateManifest(assetsFS fs.FS, cliVersion string) (*ProjectManifest, error) {
	m := &ProjectManifest{
		Version:    1,
		CLIVersion: cliVersion,
	}

	for _, mapping := range assetMappings {
		entries, err := walkAssetDir(assetsFS, mapping.sourceDir, mapping.targetDir, mapping.fileType)
		if err != nil {
			return nil, fmt.Errorf("walking %s: %w", mapping.sourceDir, err)
		}
		m.Files = append(m.Files, entries...)
	}

	agentsEntry, err := makeAgentsEntry(assetsFS)
	if err != nil {
		return nil, fmt.Errorf("reading agents.md: %w", err)
	}
	if agentsEntry != nil {
		m.Files = append(m.Files, *agentsEntry)
	}

	sort.Slice(m.Files, func(i, j int) bool {
		return m.Files[i].Source < m.Files[j].Source
	})

	return m, nil
}

func walkAssetDir(assetsFS fs.FS, sourceDir, targetDir, fileType string) ([]FileEntry, error) {
	var entries []FileEntry

	err := fs.WalkDir(assetsFS, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		content, err := fs.ReadFile(assetsFS, path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		relPath := strings.TrimPrefix(path, sourceDir+"/")
		targetPath := filepath.Join(targetDir, relPath)

		entries = append(entries, FileEntry{
			Source: path,
			Target: targetPath,
			SHA256: sha256Hex(content),
			Type:   fileType,
		})
		return nil
	})

	return entries, err
}

func makeAgentsEntry(assetsFS fs.FS) (*FileEntry, error) {
	content, err := fs.ReadFile(assetsFS, "assets/AGENTS.md")
	if err != nil {
		return nil, nil
	}

	return &FileEntry{
		Source: "assets/AGENTS.md",
		Target: "AGENTS.md",
		SHA256: sha256Hex(content),
		Type:   "agents_block",
		Markers: &BlockMarkers{
			Start: "<!-- FLOWFORGE:START -->",
			End:   "<!-- FLOWFORGE:END -->",
		},
	}, nil
}

func CompareManifests(old, new *ProjectManifest) *DiffResult {
	result := &DiffResult{}

	oldMap := make(map[string]FileEntry)
	for _, f := range old.Files {
		oldMap[f.Source] = f
	}

	newMap := make(map[string]FileEntry)
	for _, f := range new.Files {
		newMap[f.Source] = f
	}

	for _, newFile := range new.Files {
		oldFile, exists := oldMap[newFile.Source]
		if !exists {
			result.Added = append(result.Added, newFile)
			continue
		}

		if newFile.SHA256 != oldFile.SHA256 {
			if oldFile.SHA256 != newFile.SHA256 {
				oldOnDisk, err := readTargetFile(oldFile.Target)
				if err == nil && oldOnDisk != oldFile.SHA256 {
					result.Conflict = append(result.Conflict, newFile)
					continue
				}
			}
			result.Updated = append(result.Updated, newFile)
		}
	}

	for _, oldFile := range old.Files {
		if _, exists := newMap[oldFile.Source]; !exists {
			result.Removed = append(result.Removed, oldFile)
		}
	}

	return result
}

func readTargetFile(targetPath string) (string, error) {
	content, err := os.ReadFile(targetPath)
	if err != nil {
		return "", err
	}
	return sha256Hex(content), nil
}

func sha256Hex(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func (d *DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Updated) > 0 || len(d.Conflict) > 0 || len(d.Removed) > 0
}

func (d *DiffResult) Summary() string {
	var parts []string
	if len(d.Added) > 0 {
		parts = append(parts, fmt.Sprintf("%d added", len(d.Added)))
	}
	if len(d.Updated) > 0 {
		parts = append(parts, fmt.Sprintf("%d updated", len(d.Updated)))
	}
	if len(d.Conflict) > 0 {
		parts = append(parts, fmt.Sprintf("%d conflict", len(d.Conflict)))
	}
	if len(d.Removed) > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", len(d.Removed)))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}
