package command

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
	"flowforge/internal/orchestration"
	"flowforge/internal/version"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

const (
	orchestrationBlockStart = "<!-- FLOWFORGE:ORCHESTRATION:START -->"
	orchestrationBlockEnd   = "<!-- FLOWFORGE:ORCHESTRATION:END -->"
)

type hostSet map[string]bool

type syncOptions struct {
	forced  []string
	removed []string
	dryRun  bool
	adopt   bool
}

func newSyncCmd() *cobra.Command {
	var opts syncOptions
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize FlowForge project facilities",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := currentConfigService()
			if err != nil {
				return fmt.Errorf("finding project root: %w", err)
			}
			defer svc.Close()
			return syncProject(cmd, svc.ProjectRoot(), opts)
		},
	}
	cmd.Flags().StringSliceVar(&opts.forced, "host", nil, "Enable host facilities (opencode,codex)")
	cmd.Flags().StringSliceVar(&opts.removed, "without-host", nil, "Disable and remove managed host facilities")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Show changes without writing files")
	cmd.Flags().BoolVar(&opts.adopt, "adopt", false, "Replace existing FlowForge-named agent files and manage them")
	return cmd
}

func syncProject(cmd *cobra.Command, root string, opts syncOptions) error {
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		return fmt.Errorf("loading project manifest: %w", err)
	}
	disabled := makeHostSet(manifest.DisabledHosts)
	pending := makeHostSet(manifest.PendingHosts)
	for _, host := range opts.forced {
		if !validHost(host) {
			return fmt.Errorf("unsupported host %q", host)
		}
		delete(disabled, host)
		delete(pending, host)
	}
	for _, host := range opts.removed {
		if !validHost(host) {
			return fmt.Errorf("unsupported host %q", host)
		}
		disabled[host] = true
		delete(pending, host)
	}
	hosts, err := detectHosts(root, manifest)
	if err != nil {
		return err
	}
	for _, host := range opts.forced {
		hosts[host] = true
	}
	for host := range disabled {
		delete(hosts, host)
	}
	for host := range pending {
		if hostEvidenceExists(root, host) {
			delete(pending, host)
		} else {
			delete(hosts, host)
		}
	}

	if opts.dryRun {
		if err := previewAssetUpdates(cmd, root, manifest); err != nil {
			return err
		}
	} else if _, err := applyAssetUpdates(root); err != nil {
		return err
	} else if manifest, err = core.LoadProjectManifest(root); err != nil {
		return err
	}

	policy := orchestration.DefaultPolicy()
	desired := map[string]managedContent{}
	if hosts["opencode"] {
		files, renderErr := orchestration.RenderOpenCode(policy)
		if renderErr != nil {
			return renderErr
		}
		for name, content := range files {
			desired[filepath.Join(".opencode", "agents", name)] = managedContent{"generated/opencode/" + name, "opencode_agent", content}
		}
	}
	if hosts["codex"] {
		files, renderErr := orchestration.RenderCodex(policy)
		if renderErr != nil {
			return renderErr
		}
		for name, content := range files {
			desired[filepath.Join(".codex", "agents", name)] = managedContent{"generated/codex/" + name, "codex_agent", content}
		}
	}
	available, failed, err := reconcileHostFiles(cmd, root, manifest, desired, opts)
	if err != nil {
		return err
	}
	for host := range failed {
		pending[host] = true
	}
	for host := range available {
		delete(pending, host)
	}
	if err := reconcileOrchestrationBlock(cmd, root, manifest, available, opts.dryRun); err != nil {
		return err
	}
	if opts.dryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "Dry run; no files were changed.")
		return nil
	}
	manifest.DisabledHosts = sortedHosts(disabled)
	manifest.PendingHosts = sortedHosts(pending)
	manifest.CLIVersion = version.Version
	if err := manifest.Save(root); err != nil {
		return err
	}
	if len(available) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No supported agent host enabled; base facilities synchronized. Run `flowforge sync` after configuring OpenCode or Codex.")
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Synchronized hosts: %s\n", strings.Join(sortedHosts(available), ", "))
	}
	return nil
}

type managedContent struct {
	source  string
	kind    string
	content []byte
}

func validHost(host string) bool { return host == "opencode" || host == "codex" }

func makeHostSet(values []string) hostSet {
	result := hostSet{}
	for _, value := range values {
		if validHost(value) {
			result[value] = true
		}
	}
	return result
}

func detectHosts(root string, manifest *core.ProjectManifest) (hostSet, error) {
	hosts := hostSet{}
	for _, entry := range manifest.Files {
		switch entry.Type {
		case "opencode_agent":
			hosts["opencode"] = true
		case "codex_agent":
			hosts["codex"] = true
		}
	}
	checks := []struct {
		host      string
		path      string
		directory bool
	}{
		{"opencode", ".opencode", true},
		{"opencode", "opencode.json", false},
		{"opencode", "opencode.jsonc", false},
		{"codex", ".codex", true},
		{"codex", filepath.Join(".codex", "config.toml"), false},
	}
	for _, check := range checks {
		info, err := os.Stat(filepath.Join(root, check.path))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("checking %s host evidence %s: %w", check.host, check.path, err)
		}
		if check.directory != info.IsDir() {
			return nil, fmt.Errorf("invalid %s host evidence %s", check.host, check.path)
		}
		hosts[check.host] = true
	}
	return hosts, nil
}

func hostEvidenceExists(root, host string) bool {
	paths := []string{}
	if host == "opencode" {
		paths = []string{".opencode", "opencode.json", "opencode.jsonc"}
	} else if host == "codex" {
		paths = []string{".codex", filepath.Join(".codex", "config.toml")}
	}
	for _, path := range paths {
		if _, err := os.Stat(filepath.Join(root, path)); err == nil {
			return true
		}
	}
	return false
}

func sortedHosts(hosts hostSet) []string {
	result := make([]string, 0, len(hosts))
	for host := range hosts {
		result = append(result, host)
	}
	sort.Strings(result)
	return result
}

func previewAssetUpdates(cmd *cobra.Command, root string, old *core.ProjectManifest) error {
	desired, err := core.GenerateManifest(embeddedAssets, version.Version)
	if err != nil {
		return fmt.Errorf("generating manifest: %w", err)
	}
	for _, entry := range old.Files {
		if isDynamicEntry(entry) {
			desired.Files = append(desired.Files, entry)
		}
	}
	diff := core.CompareManifests(old, desired, root)
	for _, entry := range diff.Added {
		fmt.Fprintf(cmd.OutOrStdout(), "+ %s\n", entry.Target)
	}
	for _, entry := range diff.Updated {
		fmt.Fprintf(cmd.OutOrStdout(), "~ %s\n", entry.Target)
	}
	for _, entry := range diff.Conflict {
		fmt.Fprintf(cmd.ErrOrStderr(), "! conflict: %s (preserved)\n", entry.Target)
	}
	return nil
}

func isDynamicEntry(entry core.FileEntry) bool {
	return entry.Type == "opencode_agent" || entry.Type == "codex_agent" || entry.Type == "orchestration_block"
}

func reconcileHostFiles(cmd *cobra.Command, root string, manifest *core.ProjectManifest, desired map[string]managedContent, opts syncOptions) (hostSet, hostSet, error) {
	old := map[string]core.FileEntry{}
	for _, entry := range manifest.Files {
		if entry.Type == "opencode_agent" || entry.Type == "codex_agent" {
			old[entry.Target] = entry
		}
	}
	kept := make([]core.FileEntry, 0, len(manifest.Files)+len(desired))
	for _, entry := range manifest.Files {
		if entry.Type != "opencode_agent" && entry.Type != "codex_agent" {
			kept = append(kept, entry)
		}
	}
	available := hostSet{}
	failed := hostSet{}
	targets := make([]string, 0, len(desired))
	for target := range desired {
		targets = append(targets, target)
	}
	sort.Strings(targets)
	for _, target := range targets {
		item := desired[target]
		path := filepath.Join(root, target)
		data, readErr := os.ReadFile(path)
		entry, managed := old[target]
		modelPreserved := false
		if readErr == nil && managed {
			content, preserved, err := preserveAgentModel(item.kind, data, item.content)
			if err != nil {
				return nil, nil, fmt.Errorf("preserving model configuration in %s: %w", target, err)
			}
			item.content = content
			modelPreserved = preserved
		}
		adoptable := readErr == nil && (opts.adopt || knownV310Skeleton(target, data))
		if readErr == nil && !managed && !adoptable {
			fmt.Fprintf(cmd.ErrOrStderr(), "! conflict: %s (not managed; use --adopt to replace)\n", target)
			failed[hostForKind(item.kind)] = true
			continue
		}
		if readErr == nil && managed && core.SHA256Hex(data) != entry.SHA256 && !opts.adopt && !modelPreserved {
			fmt.Fprintf(cmd.ErrOrStderr(), "! conflict: %s (preserved)\n", target)
			kept = append(kept, entry)
			available[hostForKind(item.kind)] = true
			continue
		}
		if readErr != nil && !os.IsNotExist(readErr) {
			return nil, nil, fmt.Errorf("reading %s: %w", target, readErr)
		}
		if readErr != nil || core.SHA256Hex(data) != core.SHA256Hex(item.content) {
			prefix := "~"
			if readErr == nil && adoptable {
				prefix = "adopt"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", prefix, target)
			if !opts.dryRun {
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					return nil, nil, err
				}
				if err := os.WriteFile(path, item.content, 0644); err != nil {
					return nil, nil, err
				}
			}
		}
		kept = append(kept, core.FileEntry{Source: item.source, Target: target, SHA256: core.SHA256Hex(item.content), Type: item.kind})
		available[hostForKind(item.kind)] = true
	}
	for target, entry := range old {
		if _, ok := desired[target]; ok {
			continue
		}
		path := filepath.Join(root, target)
		data, readErr := os.ReadFile(path)
		if readErr != nil && !os.IsNotExist(readErr) {
			return nil, nil, fmt.Errorf("reading %s: %w", target, readErr)
		}
		if readErr == nil && core.SHA256Hex(data) != entry.SHA256 {
			fmt.Fprintf(cmd.ErrOrStderr(), "! conflict: %s (preserved)\n", target)
			kept = append(kept, entry)
			continue
		}
		if readErr == nil {
			fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", target)
			if !opts.dryRun {
				if err := os.Remove(path); err != nil {
					return nil, nil, err
				}
			}
		}
	}
	if !opts.dryRun {
		manifest.Files = kept
	}
	for host := range failed {
		delete(available, host)
	}
	return available, failed, nil
}

// preserveAgentModel carries only an explicit host model selection into the
// regenerated agent. All other generated configuration and prompt content is
// intentionally replaced by the current FlowForge version.
func preserveAgentModel(kind string, existing, generated []byte) ([]byte, bool, error) {
	switch kind {
	case "opencode_agent":
		return preserveOpenCodeModel(existing, generated)
	case "codex_agent":
		return preserveCodexModel(existing, generated)
	default:
		return generated, false, nil
	}
}

func preserveOpenCodeModel(existing, generated []byte) ([]byte, bool, error) {
	existingFrontmatter, _, ok := splitOpenCodeFrontmatter(existing)
	if !ok {
		return generated, false, nil
	}
	var oldDocument yaml.Node
	if err := yaml.Unmarshal(existingFrontmatter, &oldDocument); err != nil {
		return nil, false, fmt.Errorf("parsing existing YAML frontmatter: %w", err)
	}
	oldModel := yamlMappingValue(&oldDocument, "model")
	if oldModel == nil {
		return generated, false, nil
	}
	generatedFrontmatter, generatedBody, ok := splitOpenCodeFrontmatter(generated)
	if !ok {
		return nil, false, fmt.Errorf("generated YAML frontmatter is missing")
	}
	var newDocument yaml.Node
	if err := yaml.Unmarshal(generatedFrontmatter, &newDocument); err != nil {
		return nil, false, fmt.Errorf("parsing generated YAML frontmatter: %w", err)
	}
	newMapping := yamlDocumentMapping(&newDocument)
	if newMapping == nil {
		return nil, false, fmt.Errorf("generated YAML frontmatter is not a mapping")
	}
	model := *oldModel
	if current := yamlMappingValue(&newDocument, "model"); current != nil {
		*current = model
	} else {
		newMapping.Content = append(newMapping.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "model"}, &model)
	}
	encoded, err := yaml.Marshal(&newDocument)
	if err != nil {
		return nil, false, fmt.Errorf("encoding merged YAML frontmatter: %w", err)
	}
	return append(append([]byte("---\n"), encoded...), append([]byte("---\n"), generatedBody...)...), true, nil
}

func splitOpenCodeFrontmatter(content []byte) (frontmatter, body []byte, ok bool) {
	if !bytes.HasPrefix(content, []byte("---\n")) {
		return nil, nil, false
	}
	rest := content[len("---\n"):]
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		return nil, nil, false
	}
	return rest[:end], rest[end+len("\n---\n"):], true
}

func yamlDocumentMapping(document *yaml.Node) *yaml.Node {
	if document.Kind != yaml.DocumentNode || len(document.Content) != 1 {
		return nil
	}
	if document.Content[0].Kind != yaml.MappingNode {
		return nil
	}
	return document.Content[0]
}

func yamlMappingValue(document *yaml.Node, key string) *yaml.Node {
	mapping := yamlDocumentMapping(document)
	if mapping == nil {
		return nil
	}
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			return mapping.Content[index+1]
		}
	}
	return nil
}

func preserveCodexModel(existing, generated []byte) ([]byte, bool, error) {
	var oldDocument map[string]any
	if err := toml.Unmarshal(existing, &oldDocument); err != nil {
		return nil, false, fmt.Errorf("parsing existing TOML: %w", err)
	}
	model, hasModel := oldDocument["model"]
	if !hasModel {
		return generated, false, nil
	}
	var newDocument map[string]any
	if err := toml.Unmarshal(generated, &newDocument); err != nil {
		return nil, false, fmt.Errorf("parsing generated TOML: %w", err)
	}
	newDocument["model"] = model
	merged, err := toml.Marshal(newDocument)
	if err != nil {
		return nil, false, fmt.Errorf("encoding merged TOML: %w", err)
	}
	return merged, true, nil
}

func hostForKind(kind string) string {
	if kind == "codex_agent" {
		return "codex"
	}
	return "opencode"
}

func knownV310Skeleton(target string, data []byte) bool {
	name := filepath.Base(target)
	role, ok := map[string]struct {
		id, display, profile, skill, edit string
	}{
		"flowforge-design-analyst.md": {"design-analyst", "Design Analyst", "high-capability", "flowforge-design", "deny"},
		"flowforge-executor.md":       {"executor", "Executor", "tool-capable", "flowforge-implement", "allow"},
	}[name]
	if !ok {
		return false
	}
	permissions := "read: allow\n  edit: " + role.edit + "\n  task: deny\n  question: deny\n  skill: allow"
	body := fmt.Sprintf("---\nname: flowforge-%s\ndescription: FlowForge %s role.\nmode: subagent\npermission:\n  %s\n---\n\n# FlowForge %s\n\nActive Role: %s\nModel Profile: %s\nDefault Skill: %s\n\nRead the Proposal Journal and referenced artifacts. Follow the installed FlowForge Skill and return control to the Coordinator. Do not delegate or ask the user directly.\n", role.id, role.display, strings.ReplaceAll(permissions, "\n", "\n  "), role.display, role.display, role.profile, role.skill)
	return core.SHA256Hex(data) == core.SHA256Hex([]byte(body))
}

func reconcileOrchestrationBlock(cmd *cobra.Command, root string, manifest *core.ProjectManifest, hosts hostSet, dryRun bool) error {
	var old *core.FileEntry
	filtered := make([]core.FileEntry, 0, len(manifest.Files)+1)
	for _, entry := range manifest.Files {
		if entry.Type == "orchestration_block" {
			copy := entry
			old = &copy
			continue
		}
		filtered = append(filtered, entry)
	}
	path := filepath.Join(root, "AGENTS.md")
	file, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	current, found, extractErr := core.ExtractMarkedBlock(file, orchestrationBlockStart, orchestrationBlockEnd)
	if extractErr != nil {
		return extractErr
	}
	if old != nil && (!found || core.SHA256Hex(current) != old.SHA256) {
		fmt.Fprintln(cmd.ErrOrStderr(), "! conflict: AGENTS.md orchestration block (preserved)")
		filtered = append(filtered, *old)
		if !dryRun {
			manifest.Files = filtered
		}
		return nil
	}
	if len(hosts) == 0 {
		if found {
			fmt.Fprintln(cmd.OutOrStdout(), "- AGENTS.md orchestration block")
			if !dryRun {
				if err := core.RemoveMarkedBlock(path, orchestrationBlockStart, orchestrationBlockEnd); err != nil {
					return err
				}
			}
		}
		if !dryRun {
			manifest.Files = filtered
		}
		return nil
	}
	rules := []byte(renderOrchestrationRules(hosts))
	if !found || core.SHA256Hex(current) != core.SHA256Hex(rules) {
		fmt.Fprintln(cmd.OutOrStdout(), "~ AGENTS.md orchestration block")
		if !dryRun {
			if err := core.ApplyMarkedBlock(path, orchestrationBlockStart, orchestrationBlockEnd, rules); err != nil {
				return err
			}
		}
	}
	filtered = append(filtered, core.FileEntry{Source: "generated/AGENTS.orchestration.md", Target: "AGENTS.md", SHA256: core.SHA256Hex(rules), Type: "orchestration_block", Markers: &core.BlockMarkers{Start: orchestrationBlockStart, End: orchestrationBlockEnd}})
	if !dryRun {
		manifest.Files = filtered
	}
	return nil
}

func renderOrchestrationRules(hosts hostSet) string {
	reviewerRule := ""
	for _, role := range orchestration.DefaultPolicy().Roles {
		if role.ID == "reviewer" && role.Enabled {
			reviewerRule = "- Run `context risk-review` after implementation and use `flowforge-reviewer` when it returns `review_required`.\n"
		}
	}
	coordinatorRule := ""
	if hosts["opencode"] || hosts["codex"] {
		coordinatorRule = "- Use `flowforge-coordinator` as the primary FlowForge routing agent when your host supports selecting a project primary agent; FlowForge does not change the host default automatically.\n"
	}
	return fmt.Sprintf("## FlowForge Subagents\n\nInstalled hosts: %s\n\n"+
		coordinatorRule+
		"- The primary Coordinator is the only role that talks to the user and delegates; workers never delegate or ask the user directly.\n"+
		"- Use `flowforge-design-analyst` for Proposal design, investigation, architecture, impact analysis, and replanning.\n"+
		"- Use `flowforge-executor` only after `context preflight` returns `allow` for a planned FEATURE Step and the user explicitly requested implementation.\n"+
		reviewerRule+
		"- Run `context risk-review` after implementation; when review is required but no Reviewer is installed, the primary agent performs the read-only conformance review.\n"+
		"- Read `journal recent` before delegation. Proposal, FEATURE, Step, History, and Verification artifacts override Journal summaries.\n"+
		"- After worker completion, inspect artifact state and verification evidence, then append one concise Journal entry.\n",
		strings.Join(sortedHosts(hosts), ", "))
}
