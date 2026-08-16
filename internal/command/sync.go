package command

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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
	forced         []string
	removed        []string
	dryRun         bool
	adopt          bool
	explicitEnable bool
}

var (
	renderOpenCode = orchestration.RenderOpenCodeOutput
	renderCodex    = orchestration.RenderCodexOutput
	saveManifest   = core.SaveManifestAtomic
)

func newSyncCmd() *cobra.Command {
	var opts syncOptions
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize FlowForge project facilities",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("host") || cmd.Flags().Changed("without-host") {
				return fmt.Errorf("sync --host/--without-host are legacy flags; use `flowforge subagent enable --host <host>` or `flowforge subagent disable --host <host>`")
			}
			svc, err := currentConfigService()
			if err != nil {
				return fmt.Errorf("finding project root: %w", err)
			}
			defer svc.Close()
			return syncProject(cmd, svc.ProjectRoot(), opts)
		},
	}
	cmd.Flags().StringSliceVar(&opts.forced, "host", nil, "Legacy; use `subagent enable --host` (not accepted by sync)")
	cmd.Flags().StringSliceVar(&opts.removed, "without-host", nil, "Legacy; use `subagent disable --host` (not accepted by sync)")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Show changes without writing files")
	cmd.Flags().BoolVar(&opts.adopt, "adopt", false, "Replace existing FlowForge-named agent files and manage them")
	return cmd
}

func syncProject(cmd *cobra.Command, root string, opts syncOptions) (err error) {
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		return fmt.Errorf("loading project manifest: %w", err)
	}
	pending := makeHostSet(manifest.PendingHosts)
	for _, host := range opts.forced {
		if !validHost(host) {
			return fmt.Errorf("unsupported host %q", host)
		}
		setHostIntent(manifest, host, core.HostEnabled)
		delete(pending, host)
	}
	for _, host := range opts.removed {
		if !validHost(host) {
			return fmt.Errorf("unsupported host %q", host)
		}
		setHostIntent(manifest, host, core.HostDisabled)
		delete(pending, host)
	}
	hosts := enabledHostSet(manifest)
	if len(opts.forced) > 0 {
		// An explicit enable is scoped to exactly the hosts named by the
		// invocation. Existing intent on another host is not an implicit
		// request to render or reconcile that host.
		hosts = makeHostSet(opts.forced)
	}

	policy := orchestration.DefaultPolicy()
	desired := map[string]managedContent{}
	rendererHosts := map[string]string{}
	if hosts[core.HostOpenCode] {
		output, renderErr := renderOpenCode(policy)
		if renderErr != nil {
			return fmt.Errorf("rendering opencode: %w", renderErr)
		}
		rendererHosts[core.HostOpenCode] = output.RendererVersion + ":" + output.PolicyDigest
		for _, file := range output.Files {
			desired[filepath.Join(".opencode", "agents", file.Source)] = managedContent{"generated/opencode/" + file.Source, "opencode_agent", file.Content}
		}
	}
	if hosts[core.HostCodex] {
		output, renderErr := renderCodex(policy)
		if renderErr != nil {
			return fmt.Errorf("rendering codex: %w", renderErr)
		}
		rendererHosts[core.HostCodex] = output.RendererVersion + ":" + output.PolicyDigest
		for _, file := range output.Files {
			desired[filepath.Join(".codex", "agents", file.Source)] = managedContent{"generated/codex/" + file.Source, "codex_agent", file.Content}
		}
	}
	if opts.dryRun {
		if err := previewAssetUpdates(cmd, root, manifest); err != nil {
			return err
		}
		// Dry-run must execute the same host reconciliation planner as the real
		// path. It deliberately skips asset application and the final manifest
		// save, but still reports conflicts, additions, removals, and the
		// orchestration-block plan.
		available, failed, err := reconcileHostFiles(cmd, root, manifest, desired, opts)
		if err != nil {
			return err
		}
		if opts.explicitEnable && len(failed) > 0 {
			return fmt.Errorf("explicit enable failed for host(s): %s", strings.Join(sortedHosts(failed), ", "))
		}
		if err := reconcileOrchestrationBlock(cmd, root, manifest, available, true); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Dry run; no files were changed.")
		return nil
	}
	snapshot, err := captureSyncSnapshot(root, manifest, desired)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if committed || err == nil {
			return
		}
		if rollbackErr := snapshot.restore(); rollbackErr != nil {
			err = fmt.Errorf("%w (rollback failed: %v)", err, rollbackErr)
		}
	}()

	if report, err := applyAssetUpdates(root, opts.adopt); err != nil {
		return err
	} else if manifest, err = core.LoadProjectManifest(root); err != nil {
		return err
	} else if report != nil {
		for _, entry := range report.Conflict {
			fmt.Fprintf(cmd.ErrOrStderr(), "! conflict: %s -> %s (preserved)\n", entry.Source, entry.Target)
		}
	}
	// Asset reconciliation may reload the manifest. Reapply this invocation's
	// explicit intent after that reload; disk evidence remains irrelevant.
	for _, host := range opts.forced {
		setHostIntent(manifest, host, core.HostEnabled)
	}
	for _, host := range opts.removed {
		setHostIntent(manifest, host, core.HostDisabled)
	}

	available, failed, err := reconcileHostFiles(cmd, root, manifest, desired, opts)
	if err != nil {
		return err
	}
	if opts.explicitEnable && len(failed) > 0 {
		return fmt.Errorf("explicit enable failed for host(s): %s", strings.Join(sortedHosts(failed), ", "))
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
	manifest.PendingHosts = sortedHosts(pending)
	manifest.DisabledHosts = disabledHostList(manifest)
	manifest.CLIVersion = version.Version
	manifest.Renderer.Hosts = rendererHosts
	manifest.Renderer.PolicyDigest = core.SHA256Hex([]byte(strings.Join(sortedRendererHosts(rendererHosts), "\n")))
	if err := saveManifest(root, manifest); err != nil {
		return err
	}
	committed = true
	if len(available) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No supported agent host enabled; base facilities synchronized. Run `flowforge sync` after configuring OpenCode or Codex.")
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Synchronized hosts: %s\n", strings.Join(sortedHosts(available), ", "))
		for _, host := range sortedHosts(available) {
			fmt.Fprintf(cmd.OutOrStdout(), "Enforcement %s: %s\n", host, orchestration.EnforcementSummary(host))
		}
	}
	return nil
}

func sortedRendererHosts(hosts map[string]string) []string {
	keys := make([]string, 0, len(hosts))
	for host := range hosts {
		keys = append(keys, host+"="+hosts[host])
	}
	sort.Strings(keys)
	return keys
}

type syncFileSnapshot struct {
	path   string
	data   []byte
	exists bool
}

type syncSnapshot struct {
	files []syncFileSnapshot
}

func captureSyncSnapshot(root string, manifest *core.ProjectManifest, desired map[string]managedContent) (*syncSnapshot, error) {
	agentsPath, err := syncProjectPath(root, "AGENTS.md")
	if err != nil {
		return nil, err
	}
	paths := map[string]bool{
		filepath.Join(root, ".flowforge", core.ManifestFileName): true,
		agentsPath: true,
	}
	for _, entry := range manifest.Files {
		path, pathErr := syncProjectPath(root, entry.Target)
		if pathErr != nil {
			return nil, pathErr
		}
		paths[path] = true
	}
	generated, err := core.GenerateManifest(embeddedAssets, version.Version)
	if err != nil {
		return nil, fmt.Errorf("generating manifest snapshot: %w", err)
	}
	for _, entry := range generated.Files {
		path, pathErr := syncProjectPath(root, entry.Target)
		if pathErr != nil {
			return nil, pathErr
		}
		paths[path] = true
	}
	for target := range desired {
		path, pathErr := syncProjectPath(root, target)
		if pathErr != nil {
			return nil, pathErr
		}
		paths[path] = true
	}
	snapshot := &syncSnapshot{files: make([]syncFileSnapshot, 0, len(paths))}
	for path := range paths {
		data, readErr := os.ReadFile(path)
		if readErr == nil {
			snapshot.files = append(snapshot.files, syncFileSnapshot{path: path, data: append([]byte(nil), data...), exists: true})
			continue
		}
		if !os.IsNotExist(readErr) {
			return nil, fmt.Errorf("reading snapshot file %s: %w", path, readErr)
		}
		snapshot.files = append(snapshot.files, syncFileSnapshot{path: path})
	}
	return snapshot, nil
}

func syncProjectPath(root, target string) (string, error) {
	path, err := core.ProjectPath(root, target)
	if err != nil {
		return "", fmt.Errorf("validating sync target %q: %w", target, err)
	}
	return path, nil
}

func (s *syncSnapshot) restore() error {
	for _, file := range s.files {
		if file.exists {
			if err := os.MkdirAll(filepath.Dir(file.path), 0755); err != nil {
				return fmt.Errorf("restoring directory for %s: %w", file.path, err)
			}
			if err := os.WriteFile(file.path, file.data, 0644); err != nil {
				return fmt.Errorf("restoring %s: %w", file.path, err)
			}
			continue
		}
		if err := os.Remove(file.path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("removing created file %s: %w", file.path, err)
		}
	}
	return nil
}

func setHostIntent(manifest *core.ProjectManifest, host, intent string) {
	if host == core.HostOpenCode {
		manifest.HostIntent.OpenCode = intent
	} else if host == core.HostCodex {
		manifest.HostIntent.Codex = intent
	}
	if intent == core.HostEnabled {
		manifest.Mode = core.ManifestModeSubagent
	}
}

func enabledHostSet(manifest *core.ProjectManifest) hostSet {
	hosts := hostSet{}
	if manifest.Mode == core.ManifestModeSubagent && manifest.HostIntent.OpenCode == core.HostEnabled {
		hosts[core.HostOpenCode] = true
	}
	if manifest.Mode == core.ManifestModeSubagent && manifest.HostIntent.Codex == core.HostEnabled {
		hosts[core.HostCodex] = true
	}
	return hosts
}

func disabledHostList(manifest *core.ProjectManifest) []string {
	disabled := hostSet{}
	if manifest.HostIntent.OpenCode == core.HostDisabled {
		disabled[core.HostOpenCode] = true
	}
	if manifest.HostIntent.Codex == core.HostDisabled {
		disabled[core.HostCodex] = true
	}
	return sortedHosts(disabled)
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
	evidence, err := DetectHostEvidence(root, manifest)
	if err != nil {
		return nil, err
	}
	for _, item := range evidence {
		if item.Detected {
			hosts[item.Host] = true
		}
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
		fmt.Fprintf(cmd.ErrOrStderr(), "! conflict: %s -> %s (preserved)\n", entry.Source, entry.Target)
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
		path, pathErr := syncProjectPath(root, target)
		if pathErr != nil {
			return nil, nil, pathErr
		}
		data, readErr := os.ReadFile(path)
		entry, managed := old[target]
		modelPreserved := false
		generatedLike := readErr == nil && knownGeneratedAgent(target, data, item.content)
		if readErr == nil && (managed || generatedLike) {
			content, preserved, err := preserveAgentModel(item.kind, data, item.content)
			if err != nil {
				return nil, nil, fmt.Errorf("preserving model configuration in %s: %w", target, err)
			}
			item.content = content
			modelPreserved = preserved
		}
		adoptable := readErr == nil && (opts.adopt || knownV310Skeleton(target, data) || generatedLike)
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
		kept = append(kept, core.FileEntry{Source: item.source, Target: target, SHA256: core.SHA256Hex(item.content), Type: item.kind, Host: hostForKind(item.kind)})
		available[hostForKind(item.kind)] = true
	}
	for target, entry := range old {
		if _, ok := desired[target]; ok {
			continue
		}
		if _, err := core.ValidateManifestTarget(manifest, target); err != nil {
			return nil, nil, fmt.Errorf("planning removal of %s: %w", target, err)
		}
		path, pathErr := syncProjectPath(root, target)
		if pathErr != nil {
			return nil, nil, pathErr
		}
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
	generatedFrontmatter, generatedBody, ok := splitOpenCodeFrontmatter(generated)
	if !ok {
		return nil, false, fmt.Errorf("generated YAML frontmatter is missing")
	}
	preserved := false
	modelValue, found := frontmatterModelValue(existingFrontmatter)
	modelLine := []byte(nil)
	if found {
		modelLine = []byte("model: " + quoteYAMLString(modelValue) + "\n")
		preserved = true
	}
	lines := strings.Split(string(generatedFrontmatter), "\n")
	inserted := false
	if found {
		for index, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "model:") {
				lines[index] = strings.TrimSuffix(string(modelLine), "\n")
				inserted = true
				break
			}
		}
		if !inserted {
			for index, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "mode:") {
					lines = append(lines[:index+1], append([]string{strings.TrimSuffix(string(modelLine), "\n")}, lines[index+1:]...)...)
					inserted = true
					break
				}
			}
		}
		if !inserted {
			lines = append(lines, strings.TrimSuffix(string(modelLine), "\n"))
		}
	}
	mergedFrontmatter := []byte(strings.Join(lines, "\n"))
	merged, permissionPreserved := preserveOpenCodePermission(existing, mergedFrontmatter, generatedBody)
	return merged, preserved || permissionPreserved, nil
}

func preserveOpenCodePermission(existing, generatedFrontmatter, generatedBody []byte) ([]byte, bool) {
	var oldDocument, newDocument map[string]any
	oldFrontmatter, _, ok := splitOpenCodeFrontmatter(existing)
	if !ok || yaml.Unmarshal(oldFrontmatter, &oldDocument) != nil || yaml.Unmarshal(generatedFrontmatter, &newDocument) != nil {
		return append(append([]byte("---\n"), generatedFrontmatter...), append([]byte("\n---\n"), generatedBody...)...), false
	}
	oldPermission, found := oldDocument["permission"]
	if !found || reflect.DeepEqual(oldPermission, newDocument["permission"]) || isLegacyOpenCodePermission(existing, oldPermission) {
		return append(append([]byte("---\n"), generatedFrontmatter...), append([]byte("\n---\n"), generatedBody...)...), false
	}
	newDocument["permission"] = oldPermission
	encoded, err := yaml.Marshal(newDocument)
	if err != nil {
		return append(append([]byte("---\n"), generatedFrontmatter...), append([]byte("\n---\n"), generatedBody...)...), false
	}
	return append(append([]byte("---\n"), bytes.TrimSuffix(encoded, []byte("\n"))...), append([]byte("\n---\n"), generatedBody...)...), true
}

func isLegacyOpenCodePermission(content []byte, permission any) bool {
	role := ""
	for _, candidate := range []string{"Coordinator", "Design Analyst", "Executor"} {
		if bytes.Contains(content, []byte("Active Role: "+candidate)) {
			role = candidate
			break
		}
	}
	legacy := map[string]string{
		"Coordinator": `edit: deny
task:
  "*": deny
  flowforge-design-analyst: allow
  flowforge-executor: allow
question: allow
skill: allow
`,
		"Design Analyst": `edit: allow
task: deny
question: deny
skill: allow
bash:
  "*": deny
  "git status*": allow
  "git diff*": allow
  "git log*": allow
  "git show*": allow
`,
		"Executor": `edit: allow
task: deny
question: deny
skill: allow
bash:
  "*": ask
  "git status*": allow
  "git diff*": allow
`,
	}[role]
	if legacy == "" {
		return false
	}
	var expected any
	if yaml.Unmarshal([]byte(legacy), &expected) != nil {
		return false
	}
	return reflect.DeepEqual(permission, expected)
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

func frontmatterModelValue(frontmatter []byte) (string, bool) {
	for _, line := range strings.Split(string(frontmatter), "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "model:") {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(trimmed, "model:"))
		value = strings.Trim(value, "\"'")
		return value, value != ""
	}
	return "", false
}

func quoteYAMLString(value string) string {
	return fmt.Sprintf("%q", value)
}

func preserveCodexModel(existing, generated []byte) ([]byte, bool, error) {
	var oldDocument map[string]any
	if err := toml.Unmarshal(existing, &oldDocument); err != nil {
		return nil, false, fmt.Errorf("parsing existing TOML: %w", err)
	}
	var newDocument map[string]any
	if err := toml.Unmarshal(generated, &newDocument); err != nil {
		return nil, false, fmt.Errorf("parsing generated TOML: %w", err)
	}
	preserved := false
	if model, hasModel := oldDocument["model"]; hasModel {
		newDocument["model"] = model
		preserved = true
	}
	role := fmt.Sprint(oldDocument["name"])
	legacySandbox := map[string]string{
		"flowforge-coordinator":    "read-only",
		"flowforge-design-analyst": "read-only",
		"flowforge-executor":       "workspace-write",
	}[role]
	if sandbox, found := oldDocument["sandbox_mode"]; found && sandbox != newDocument["sandbox_mode"] && fmt.Sprint(sandbox) != legacySandbox {
		newDocument["sandbox_mode"] = sandbox
		preserved = true
	}
	legacyEffort := map[string]string{
		"flowforge-coordinator":    "medium",
		"flowforge-design-analyst": "high",
		"flowforge-executor":       "medium",
	}[role]
	if effort, found := oldDocument["model_reasoning_effort"]; found && effort != newDocument["model_reasoning_effort"] && fmt.Sprint(effort) != legacyEffort {
		newDocument["model_reasoning_effort"] = effort
		preserved = true
	}
	for _, key := range []string{"approval_policy", "network_access"} {
		if value, found := oldDocument[key]; found {
			newDocument[key] = value
			preserved = true
		}
	}
	if !preserved {
		return generated, false, nil
	}
	merged, err := toml.Marshal(newDocument)
	if err != nil {
		return nil, false, fmt.Errorf("encoding merged TOML: %w", err)
	}
	return merged, true, nil
}

func knownGeneratedAgent(target string, existing, generated []byte) bool {
	name := filepath.Base(target)
	if strings.HasPrefix(name, "flowforge-") && strings.Contains(string(existing), "# FlowForge ") &&
		strings.Contains(string(existing), "Active Role:") && strings.Contains(string(existing), "## Shared Workflow") {
		return true
	}
	if filepath.Ext(target) == ".md" {
		return strings.TrimSpace(string(removeModelLine(existing))) == strings.TrimSpace(string(removeModelLine(generated)))
	}
	if filepath.Ext(target) == ".toml" {
		return strings.TrimSpace(string(removeModelLine(existing))) == strings.TrimSpace(string(removeModelLine(generated)))
	}
	return false
}

func removeModelLine(content []byte) []byte {
	lines := strings.Split(string(content), "\n")
	filtered := lines[:0]
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "model:") || strings.HasPrefix(trimmed, "model =") {
			continue
		}
		filtered = append(filtered, line)
	}
	return []byte(strings.Join(filtered, "\n"))
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
	path, pathErr := syncProjectPath(root, "AGENTS.md")
	if pathErr != nil {
		return pathErr
	}
	file, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	current, found, extractErr := core.ExtractMarkedBlock(file, orchestrationBlockStart, orchestrationBlockEnd)
	if extractErr != nil {
		return extractErr
	}
	orchestrationContent := []byte(nil)
	if found {
		orchestrationContent = current
	}
	if old != nil && (!found || core.SHA256Hex(current) != old.SHA256) {
		fmt.Fprintln(cmd.ErrOrStderr(), "! conflict: AGENTS.md orchestration block (preserved)")
	}
	if len(hosts) > 0 && len(orchestrationContent) == 0 {
		orchestrationContent = []byte(renderOrchestrationRules(hosts))
	}
	base, baseFound, err := core.ExtractMarkedBlock(file, "<!-- FLOWFORGE:START -->", "<!-- FLOWFORGE:END -->")
	if err != nil {
		return err
	}
	if len(hosts) == 0 {
		orchestrationContent = nil
	}
	merged := append([]byte(nil), base...)
	if baseFound {
		merged = removeNestedOrchestration(merged)
		if len(orchestrationContent) > 0 {
			merged = bytes.TrimRight(merged, "\n")
			if len(merged) > 0 {
				merged = append(merged, []byte("\n\n")...)
			}
			merged = append(merged, []byte(orchestrationBlockStart+"\n")...)
			merged = append(merged, orchestrationContent...)
			if merged[len(merged)-1] != '\n' {
				merged = append(merged, '\n')
			}
			merged = append(merged, []byte(orchestrationBlockEnd+"\n")...)
		}
	}
	if !baseFound {
		if len(hosts) > 0 {
			merged = append([]byte("\n"+orchestrationBlockStart+"\n"), orchestrationContent...)
			merged = append(merged, []byte(orchestrationBlockEnd+"\n")...)
		}
	}
	if !baseFound && len(hosts) == 0 {
		if !dryRun {
			manifest.Files = filtered
		}
		return nil
	}
	if !bytes.Equal(file, merged) {
		fmt.Fprintln(cmd.OutOrStdout(), "~ AGENTS.md")
	}
	if !dryRun {
		if found {
			if err := core.RemoveMarkedBlock(path, orchestrationBlockStart, orchestrationBlockEnd); err != nil {
				return err
			}
		}
		if err := core.ApplyMarkedBlock(path, "<!-- FLOWFORGE:START -->", "<!-- FLOWFORGE:END -->", merged); err != nil {
			return err
		}
	}
	if !dryRun {
		manifest.Files = filtered
	}
	if len(hosts) > 0 {
		filtered = append(filtered, core.FileEntry{Source: "generated/AGENTS.orchestration.md", Target: "AGENTS.md", SHA256: core.SHA256Hex(orchestrationContent), Type: "orchestration_block", Markers: &core.BlockMarkers{Start: orchestrationBlockStart, End: orchestrationBlockEnd}})
		if !dryRun {
			manifest.Files = filtered
		}
	}
	return nil
}

func removeNestedOrchestration(content []byte) []byte {
	start := bytes.Index(content, []byte(orchestrationBlockStart))
	end := bytes.Index(content, []byte(orchestrationBlockEnd))
	if start < 0 || end <= start {
		return content
	}
	end += len(orchestrationBlockEnd)
	if end < len(content) && content[end] == '\n' {
		end++
	}
	return append(append([]byte(nil), content[:start]...), content[end:]...)
}

func renderOrchestrationRules(hosts hostSet) string {
	reviewerRule := ""
	for _, role := range orchestration.DefaultPolicy().Roles {
		if role.ID == "reviewer" && role.Enabled {
			reviewerRule = "- Run `context risk-review` after implementation and use `flowforge-reviewer` when it returns `review_required`.\n"
		}
	}
	coordinatorRule := ""
	if hosts["opencode"] {
		coordinatorRule += "- In OpenCode, use `flowforge-coordinator` as the primary FlowForge routing agent.\n"
	}
	if hosts["codex"] {
		coordinatorRule += "- In Codex, the current main thread IS the FlowForge Coordinator by default. The installed `flowforge-coordinator` custom agent remains available as a manual fallback: do not spawn it automatically, but use it when the user explicitly requests `flowforge-coordinator`.\n" +
			"- For FlowForge-managed work in Codex, delegation is mandatory even for one Step, non-parallel work, or work the main thread could perform locally. The main thread MUST use native subagents and MUST NOT edit worker-owned design artifacts or product code.\n" +
			"- In Codex, design, decomposition, architecture, impact analysis, evidence synthesis, and replanning MUST spawn `flowforge-design-analyst`; each ready registered investigation brief MUST spawn `flowforge-investigator`; implementation MUST run `context preflight` and, only when it returns `allow` with `Handoff: required`, MUST spawn `flowforge-executor` with the exact Step context.\n" +
			"- If native Codex subagents are unavailable, return `BLOCKED` instead of silently executing worker work in the main thread.\n"
	}
	return fmt.Sprintf("## FlowForge Subagents\n\nInstalled hosts: %s\n\n"+
		coordinatorRule+
		"- The primary Coordinator is the only role that talks to the user and delegates. It is an execution scheduler: read structured analysis revision/readiness/re-entry state, show the user each background action, and dispatch only work already registered by the Design Analyst.\n"+
		"- Delegation depth is one: every worker is dispatched directly by the Coordinator; workers never delegate or ask the user directly.\n"+
		"- Use `flowforge-design-analyst` for framing, FEATURE decomposition, investigation planning, evidence synthesis, architecture, impact analysis, and replanning.\n"+
		"- Use `flowforge-investigator` only for a ready registered investigation brief; it writes only the assigned FIND and returns structured blocked, inconclusive, conflict, or decision status.\n"+
		"- Use `flowforge-executor` only after `context preflight` returns `allow` for a planned FEATURE Step and the user explicitly requested implementation.\n"+
		reviewerRule+
		"- Run `context risk-review` after implementation; when review is required but no Reviewer is installed, the primary agent performs the read-only conformance review.\n"+
		"- Read `journal recent` before delegation. Proposal, FEATURE, DEC, FIND, Step, History, and Verification own durable facts; Journal owns analysis scheduling state and artifact references.\n"+
		"- External sources require explicit authorization in the work-item brief; unavailable required access returns BLOCKED.\n"+
		"- After worker completion, inspect artifact state and verification evidence, then append one concise Journal entry.\n",
		strings.Join(sortedHosts(hosts), ", "))
}
