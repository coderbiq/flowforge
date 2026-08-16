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

const (
	ManifestVersionV1       = 1
	ManifestVersionV2       = 2
	ManifestModeNonSubagent = "non_subagent"
	ManifestModeSubagent    = "subagent"
	HostOpenCode            = "opencode"
	HostCodex               = "codex"
	HostEnabled             = "enabled"
	HostDisabled            = "disabled"
)

type FileEntry struct {
	Source  string        `yaml:"source"`
	Target  string        `yaml:"target"`
	SHA256  string        `yaml:"sha256"`
	Type    string        `yaml:"type"`
	Host    string        `yaml:"host,omitempty"`
	Dormant bool          `yaml:"dormant,omitempty"`
	Markers *BlockMarkers `yaml:"markers,omitempty"`
}

type BlockMarkers struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

type ProjectManifest struct {
	Version       int              `yaml:"version"`
	CLIVersion    string           `yaml:"cli_version"`
	DisabledHosts []string         `yaml:"disabled_hosts,omitempty"`
	PendingHosts  []string         `yaml:"pending_hosts,omitempty"`
	Mode          string           `yaml:"mode,omitempty"`
	HostIntent    HostIntent       `yaml:"host_intent,omitempty"`
	Renderer      RendererMetadata `yaml:"renderer,omitempty"`
	Files         []FileEntry      `yaml:"files"`
}

type HostIntent struct {
	OpenCode string `yaml:"opencode"`
	Codex    string `yaml:"codex"`
}

type RendererMetadata struct {
	PolicyDigest string            `yaml:"policy_digest,omitempty"`
	Hosts        map[string]string `yaml:"hosts,omitempty"`
}

type DiffResult struct {
	Added    []FileEntry
	Updated  []FileEntry
	Conflict []FileEntry
	Removed  []FileEntry
}

// MigrationResult describes a validated in-memory manifest migration. The
// caller decides when the result becomes durable.
type MigrationResult struct {
	Manifest *ProjectManifest
	Migrated bool
}

// DynamicEntriesForHost returns only registered, host-specific entries. The
// manifest is the source of truth for these entries; filesystem observations
// are deliberately not consulted.
func (m *ProjectManifest) DynamicEntriesForHost(host string) []FileEntry {
	if m == nil || (host != HostOpenCode && host != HostCodex) {
		return []FileEntry{}
	}
	entries := make([]FileEntry, 0)
	for _, entry := range m.Files {
		if !isDynamicFileEntry(entry) || entry.Host != host {
			continue
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Target == entries[j].Target {
			return entries[i].Source < entries[j].Source
		}
		return entries[i].Target < entries[j].Target
	})
	return entries
}

// DesiredHostSet contains the registered dynamic entries only when the host
// has an explicit enabled intent. non_subagent and disabled hosts are empty.
func (m *ProjectManifest) DesiredHostSet(host string) []FileEntry {
	if m == nil || m.Mode == ManifestModeNonSubagent || m.hostIntent(host) != HostEnabled {
		return []FileEntry{}
	}
	return m.DynamicEntriesForHost(host)
}

// ValidateManifestTarget verifies that target is a uniquely registered
// manifest target and returns its registration. It never treats a disk path
// as authorization to manage or delete a file.
func ValidateManifestTarget(m *ProjectManifest, target string) (FileEntry, error) {
	if m == nil {
		return FileEntry{}, fmt.Errorf("manifest is nil")
	}
	if err := validateManifestPath("target", target); err != nil {
		return FileEntry{}, err
	}
	var found *FileEntry
	for i := range m.Files {
		if m.Files[i].Target != target || !isDynamicFileEntry(m.Files[i]) {
			continue
		}
		if found != nil {
			return FileEntry{}, fmt.Errorf("duplicate dynamic target %q", target)
		}
		found = &m.Files[i]
	}
	if found == nil {
		return FileEntry{}, fmt.Errorf("target %q is not registered", target)
	}
	return *found, nil
}

// ProjectPath resolves a manifest target while keeping it inside projectRoot.
// Lexical validation alone is insufficient when a project contains symlinks.
func ProjectPath(projectRoot, target string) (string, error) {
	if err := validateManifestPath("target", target); err != nil {
		return "", err
	}
	root, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("resolving project root: %w", err)
	}
	path := filepath.Join(root, filepath.FromSlash(target))
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolving project root symlinks: %w", err)
	}
	check := path
	missing := make([]string, 0)
	for current := path; ; current = filepath.Dir(current) {
		realPath, evalErr := filepath.EvalSymlinks(current)
		if evalErr == nil {
			for i := len(missing) - 1; i >= 0; i-- {
				realPath = filepath.Join(realPath, missing[i])
			}
			check = realPath
			break
		}
		if !os.IsNotExist(evalErr) {
			return "", fmt.Errorf("resolving target %q: %w", target, evalErr)
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("resolving target parent %q: %w", target, evalErr)
		}
		missing = append(missing, filepath.Base(current))
	}
	rel, err := filepath.Rel(realRoot, check)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("target %q escapes project root", target)
	}
	return path, nil
}

type ReconcileActionKind string

const (
	ReconcileAdd    ReconcileActionKind = "add"
	ReconcileUpdate ReconcileActionKind = "update"
	ReconcileRemove ReconcileActionKind = "remove"
)

// ReconcileAction is a deterministic plan item. Remove actions can only be
// constructed from a manifest registration.
type ReconcileAction struct {
	Kind   ReconcileActionKind
	Target string
	Entry  FileEntry
}

func isDynamicFileEntry(entry FileEntry) bool {
	return entry.Type == "opencode_agent" || entry.Type == "codex_agent" || entry.Type == "orchestration_block"
}

func (m *ProjectManifest) hostIntent(host string) string {
	switch host {
	case HostOpenCode:
		return m.HostIntent.OpenCode
	case HostCodex:
		return m.HostIntent.Codex
	default:
		return HostDisabled
	}
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

	if m.Version == ManifestVersionV1 {
		return MigrateManifestV1(&m)
	}
	if err := m.NormalizeAndValidate(); err != nil {
		return nil, fmt.Errorf("validating manifest: %w", err)
	}
	return &m, nil
}

// NormalizeManifest validates and returns a normalized copy without writing it.
func NormalizeManifest(m *ProjectManifest) (*ProjectManifest, error) {
	if m == nil {
		return nil, fmt.Errorf("manifest is nil")
	}
	normalized := *m
	if err := normalized.NormalizeAndValidate(); err != nil {
		return nil, err
	}
	return &normalized, nil
}

// LoadProjectManifestMigration loads and fully validates a manifest before
// returning a migration result. It never writes project files.
func LoadProjectManifestMigration(projectRoot string) (*MigrationResult, error) {
	path := filepath.Join(projectRoot, configDirName, ManifestFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest: %w", err)
	}
	var m ProjectManifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest: %w", err)
	}
	if m.Version == ManifestVersionV1 {
		migrated, err := MigrateManifestV1(&m)
		if err != nil {
			return nil, err
		}
		return &MigrationResult{Manifest: migrated, Migrated: true}, nil
	}
	normalized, err := NormalizeManifest(&m)
	if err != nil {
		return nil, fmt.Errorf("validating manifest: %w", err)
	}
	return &MigrationResult{Manifest: normalized}, nil
}

// MigrateManifestV1 converts a parsed v1 manifest in memory. It never writes
// files; callers must explicitly call Save after reviewing the result.
func MigrateManifestV1(v1 *ProjectManifest) (*ProjectManifest, error) {
	if v1 == nil {
		return nil, fmt.Errorf("manifest is nil")
	}
	m := *v1
	m.Version = ManifestVersionV2
	m.Mode = ManifestModeNonSubagent
	m.HostIntent = HostIntent{OpenCode: HostDisabled, Codex: HostDisabled}
	m.Renderer = RendererMetadata{}
	m.DisabledHosts = nil
	m.PendingHosts = nil
	for i := range m.Files {
		entry := &m.Files[i]
		if entry.Host == "" {
			entry.Host = hostForEntry(entry)
		}
		if entry.Host == "" {
			entry.Host = hostForEntry(entry)
		}
		if entry.Host != "" {
			entry.Dormant = true
		}
	}
	if err := m.NormalizeAndValidate(); err != nil {
		return nil, fmt.Errorf("migrating v1 manifest: %w", err)
	}
	return &m, nil
}

func (m *ProjectManifest) Save(projectRoot string) error {
	return SaveManifestAtomic(projectRoot, m)
}

// SaveManifestAtomic validates and atomically replaces the project manifest.
func SaveManifestAtomic(projectRoot string, m *ProjectManifest) error {
	if m == nil {
		return fmt.Errorf("manifest is nil")
	}
	normalized, err := NormalizeManifest(m)
	if err != nil {
		return fmt.Errorf("validating manifest before save: %w", err)
	}
	path := filepath.Join(projectRoot, configDirName, ManifestFileName)
	data, err := yaml.Marshal(normalized)
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating manifest dir: %w", err)
	}

	tmp, err := os.CreateTemp(dir, ".manifest-*.tmp")
	if err != nil {
		return fmt.Errorf("creating manifest temp file: %w", err)
	}
	tmpName := tmp.Name()
	if err := tmp.Chmod(0644); err != nil {
		if closeErr := tmp.Close(); closeErr != nil {
			return fmt.Errorf("setting manifest temp permissions: %w (closing: %v)", err, closeErr)
		}
		return fmt.Errorf("setting manifest temp permissions: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		if closeErr := tmp.Close(); closeErr != nil {
			return fmt.Errorf("writing manifest temp: %w (closing: %v)", err, closeErr)
		}
		return fmt.Errorf("writing manifest temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing manifest temp: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		if removeErr := os.Remove(tmpName); removeErr != nil && !os.IsNotExist(removeErr) {
			return fmt.Errorf("writing manifest: %w (cleaning temp: %v)", err, removeErr)
		}
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
		Version:    ManifestVersionV2,
		CLIVersion: cliVersion,
		Mode:       ManifestModeNonSubagent,
		HostIntent: HostIntent{OpenCode: HostDisabled, Codex: HostDisabled},
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

func (m *ProjectManifest) NormalizeAndValidate() error {
	if m.Version != ManifestVersionV2 {
		return fmt.Errorf("unsupported manifest schema version %d", m.Version)
	}
	if m.Mode != ManifestModeNonSubagent && m.Mode != ManifestModeSubagent {
		return fmt.Errorf("invalid mode %q", m.Mode)
	}
	if err := validateHostIntent(m); err != nil {
		return err
	}
	seenSources := make(map[string]struct{}, len(m.Files))
	seenTargets := make(map[string]string, len(m.Files))
	for i := range m.Files {
		entry := &m.Files[i]
		if entry.Host == "" {
			entry.Host = hostForEntry(entry)
		}
		if err := validateManifestPath("source", entry.Source); err != nil {
			return fmt.Errorf("file %d: %w", i, err)
		}
		if err := validateManifestPath("target", entry.Target); err != nil {
			return fmt.Errorf("file %d: %w", i, err)
		}
		if _, ok := seenSources[entry.Source]; ok {
			return fmt.Errorf("duplicate source %q", entry.Source)
		}
		if previous, ok := seenTargets[entry.Target]; ok {
			// AGENTS.md intentionally has separate, independently managed blocks.
			// Other target collisions are ambiguous and remain invalid.
			if entry.Target != "AGENTS.md" ||
				!((previous == "agents_block" && entry.Type == "orchestration_block") ||
					(previous == "orchestration_block" && entry.Type == "agents_block")) {
				return fmt.Errorf("duplicate target %q", entry.Target)
			}
		}
		seenSources[entry.Source] = struct{}{}
		seenTargets[entry.Target] = entry.Type
		if entry.Target == "AGENTS.md" && entry.Type == "orchestration_block" {
			continue
		}
		if entry.Host != "" && entry.Host != HostOpenCode && entry.Host != HostCodex {
			return fmt.Errorf("file %q has unknown host %q", entry.Source, entry.Host)
		}
		if entry.Markers != nil && (entry.Markers.Start == "" || entry.Markers.End == "" || entry.Markers.Start == entry.Markers.End) {
			return fmt.Errorf("file %q has invalid markers", entry.Source)
		}
		if entry.Type == "opencode_agent" && entry.Host != HostOpenCode {
			return fmt.Errorf("file %q must belong to opencode", entry.Source)
		}
		if entry.Type == "codex_agent" && entry.Host != HostCodex {
			return fmt.Errorf("file %q must belong to codex", entry.Source)
		}
	}
	if m.Mode == ManifestModeNonSubagent && (m.HostIntent.OpenCode != HostDisabled || m.HostIntent.Codex != HostDisabled) {
		return fmt.Errorf("non_subagent mode requires disabled host intents")
	}
	sort.Slice(m.Files, func(i, j int) bool {
		if m.Files[i].Source == m.Files[j].Source {
			return m.Files[i].Target < m.Files[j].Target
		}
		return m.Files[i].Source < m.Files[j].Source
	})
	return nil
}

func validateHostIntent(m *ProjectManifest) error {
	for host, intent := range map[string]string{HostOpenCode: m.HostIntent.OpenCode, HostCodex: m.HostIntent.Codex} {
		if intent != HostEnabled && intent != HostDisabled {
			return fmt.Errorf("invalid %s host intent %q", host, intent)
		}
	}
	return nil
}

func validateManifestPath(kind, value string) error {
	if value == "" || filepath.IsAbs(value) || filepath.Clean(value) != value || value == "." || strings.HasPrefix(value, ".."+string(filepath.Separator)) || strings.HasPrefix(value, "../") {
		return fmt.Errorf("invalid %s path %q", kind, value)
	}
	return nil
}

func hostForEntry(entry *FileEntry) string {
	switch entry.Type {
	case "opencode_agent":
		return HostOpenCode
	case "codex_agent":
		return HostCodex
	case "orchestration_block":
		if strings.Contains(entry.Source, "opencode") || strings.Contains(entry.Target, ".opencode") {
			return HostOpenCode
		}
		if strings.Contains(entry.Source, "codex") || strings.Contains(entry.Target, ".codex") {
			return HostCodex
		}
	}
	return ""
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
		SHA256: sha256Hex(StripBlockMarkers(content)),
		Type:   "agents_block",
		Markers: &BlockMarkers{
			Start: "<!-- FLOWFORGE:START -->",
			End:   "<!-- FLOWFORGE:END -->",
		},
	}, nil
}

func CompareManifests(old, new *ProjectManifest, projectRoot string) *DiffResult {
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
			// A newly shipped static asset must not overwrite an existing file
			// that is absent from the trusted baseline. Managed blocks are
			// reconciled separately and generated host files have their own
			// adoption rules.
			if newFile.Type != "agents_block" && newFile.Type != "opencode_agent" && newFile.Type != "codex_agent" {
				if _, err := os.Stat(filepath.Join(projectRoot, newFile.Target)); err == nil {
					result.Conflict = append(result.Conflict, newFile)
					continue
				} else if !os.IsNotExist(err) {
					result.Conflict = append(result.Conflict, newFile)
					continue
				}
			}
			result.Added = append(result.Added, newFile)
			continue
		}

		if newFile.SHA256 != oldFile.SHA256 {
			if oldFile.SHA256 != newFile.SHA256 {
				matches, err := targetMatchesEntry(oldFile, projectRoot)
				if err != nil && !os.IsNotExist(err) {
					result.Conflict = append(result.Conflict, newFile)
					continue
				}
				if err == nil && !matches {
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

func targetMatchesEntry(entry FileEntry, projectRoot string) (bool, error) {
	content, err := os.ReadFile(filepath.Join(projectRoot, entry.Target))
	if err != nil {
		return false, err
	}
	if entry.Markers != nil {
		block, found, err := ExtractMarkedBlock(content, entry.Markers.Start, entry.Markers.End)
		if err != nil {
			return false, err
		}
		if !found {
			return false, fmt.Errorf("managed block not found")
		}
		if sha256Hex(block) == entry.SHA256 {
			return true, nil
		}
		legacy := []byte(entry.Markers.Start + "\n")
		legacy = append(legacy, block...)
		legacy = append(legacy, []byte(entry.Markers.End+"\n")...)
		return sha256Hex(legacy) == entry.SHA256, nil
	}
	return sha256Hex(content) == entry.SHA256, nil
}

func sha256Hex(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256Hex returns the manifest checksum for generated managed assets.
func SHA256Hex(data []byte) string {
	return sha256Hex(data)
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
