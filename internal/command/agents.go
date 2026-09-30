package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"flowforge/internal/config"
	"flowforge/internal/subagent"
)

func newAgentsCmd() *cobra.Command {
	agents := &cobra.Command{
		Use:   "agents",
		Short: "Manage subagent definitions for Claude Code, OpenCode, Codex, and PI",
	}

	deploy := &cobra.Command{
		Use:   "deploy [name]",
		Short: "Deploy subagent definitions to host-specific directories",
		Long: `Deploy subagent definitions to .claude/agents/, .opencode/agent/, .codex/agents/, and .pi/agents/.
If [name] is specified, deploys only that subagent. Otherwise, deploys all non-disabled subagents.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return fmt.Errorf("locating project root: %w", err)
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			var targetName string
			if len(args) == 1 {
				targetName = args[0]
			}

			deployed, err := deploySubagents(projectRoot, cfg, targetName)
			if err != nil {
				return err
			}

			if len(deployed) == 0 {
				cmd.Println("No subagents deployed (all disabled or name not found)")
				return nil
			}

			hosts, err := resolveHostTargets(cfg)
			if err != nil {
				return err
			}
			cmd.Printf("✓ Deployed %d subagent(s) to %s\n", len(deployed), describeHostDirs(hosts))
			for _, name := range deployed {
				cmd.Printf("  - %s\n", name)
			}
			return nil
		},
	}

	remove := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a subagent from all host directories",
		Long: `Remove a subagent from .claude/agents/, .opencode/agent/, .codex/agents/, and .pi/agents/.
For built-in subagents, marks them as disabled in config to prevent redeployment.
For custom subagents, deletes the source definition from .flowforge/subagents/.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return fmt.Errorf("locating project root: %w", err)
			}

			name := args[0]
			isBuiltin, removedPaths, err := removeSubagent(projectRoot, name)
			if err != nil {
				return err
			}

			cmd.Printf("✓ Removed subagent %q\n", name)
			for _, path := range removedPaths {
				cmd.Printf("  - %s\n", path)
			}

			if isBuiltin {
				cmd.Printf("\nBuilt-in subagent disabled. Future 'flowforge agents deploy' will skip %q.\n", name)
			}

			return nil
		},
	}

	agents.AddCommand(deploy, remove)

	// Add status subcommand
	statusJSON := false
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Compare deployed subagent files against expected compiled content",
		Long: `Compare deployed subagent files in .claude/agents/, .opencode/agent/, .codex/agents/, and .pi/agents/
against the expected compiled content. Reports current/missing/drifted/project-owned states.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return fmt.Errorf("locating project root: %w", err)
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			result, err := computeSubagentStatus(projectRoot, cfg)
			if err != nil {
				return err
			}

			if statusJSON {
				data, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return err
				}
				cmd.Println(string(data))
			} else {
				for _, entry := range result.Entries {
					cmd.Printf("%s %s\n", entry.State, entry.Target)
				}
				if result.Current {
					cmd.Println("\n✓ All subagents current.")
				} else {
					cmd.Println("\n⚠ Some subagents missing or drifted.")
				}
			}

			if !result.Current {
				return errPolicyViolation
			}
			return nil
		},
	}
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "Output JSON")
	agents.AddCommand(statusCmd)

	return agents
}

// hostTarget describes one deployable agent host: its config key, agent
// directory relative to the project root, compiled file extension, and the
// compiler that produces host-native content. Compilers receive resolved
// per-definition options; hosts that cannot express an option ignore it.
type hostTarget struct {
	key     string
	relDir  string
	ext     string
	compile func(def *subagent.Definition, opts subagent.CompileOptions) ([]byte, error)
}

func allHostTargets() []hostTarget {
	return []hostTarget{
		{"claude", filepath.Join(".claude", "agents"), ".md", func(def *subagent.Definition, opts subagent.CompileOptions) ([]byte, error) {
			return subagent.CompileClaudeCodeWithOptions(def, opts)
		}},
		{"opencode", filepath.Join(".opencode", "agent"), ".md", func(def *subagent.Definition, opts subagent.CompileOptions) ([]byte, error) {
			return subagent.CompileOpenCodeWithOptions(def, opts)
		}},
		{"codex", filepath.Join(".codex", "agents"), ".toml", func(def *subagent.Definition, opts subagent.CompileOptions) ([]byte, error) {
			return subagent.CompileCodexWithOptions(def, opts)
		}},
		{"pi", filepath.Join(".pi", "agents"), ".md", func(def *subagent.Definition, opts subagent.CompileOptions) ([]byte, error) {
			return subagent.CompilePiWithOptions(def, opts)
		}},
	}
}

// resolveHostTargets returns the enabled host targets for a config. A nil
// agents.hosts entry means "all hosts" (backward compatibility); an explicit
// empty list or an unknown host name is a configuration error.
func resolveHostTargets(cfg *config.Config) ([]hostTarget, error) {
	all := allHostTargets()
	if cfg.Agents.Hosts == nil {
		return all, nil
	}
	if len(cfg.Agents.Hosts) == 0 {
		return nil, fmt.Errorf("agents.hosts must name at least one of claude, opencode, codex, pi")
	}
	byKey := make(map[string]hostTarget, len(all))
	for _, h := range all {
		byKey[h.key] = h
	}
	seen := make(map[string]bool, len(cfg.Agents.Hosts))
	var selected []hostTarget
	for _, name := range cfg.Agents.Hosts {
		if seen[name] {
			continue
		}
		h, ok := byKey[name]
		if !ok {
			return nil, fmt.Errorf("agents.hosts: unknown host %q (supported: claude, opencode, codex, pi)", name)
		}
		seen[name] = true
		selected = append(selected, h)
	}
	return selected, nil
}

// describeHostDirs renders the enabled host directories for user messages.
func describeHostDirs(hosts []hostTarget) string {
	parts := make([]string, len(hosts))
	for i, h := range hosts {
		parts[i] = h.relDir + string(filepath.Separator)
	}
	return strings.Join(parts, ", ")
}

// defaultTestFileGlobs is the default write-protection set for the
// flowforge-implementer agent (test author / implementer separation).
var defaultTestFileGlobs = []string{
	"**/*_test.go",
	"**/src/test/**",
	"**/src/integrationTest/**",
	"**/__tests__/**",
	"**/*.test.ts",
	"**/*.test.tsx",
	"**/*.spec.ts",
}

// validModelProfileKeys are the agents.models config keys the CLI accepts.
var validModelProfileKeys = map[string]bool{
	"tool-capable":           true,
	"tool-capable-read-only": true,
}

// defaultImplementerMaxSteps is the execution budget compiled into the
// flowforge-implementer OpenCode agent when agents.max_steps is unconfigured.
// Rationale (executor-loop-hardening design d-execution-budget): an incident
// outlier session burned 893 tool calls while normal sessions stay in the
// 51-149 range; 200 covers the normal upper bound and truncates the outlier
// tail via the host's native `steps` limit.
const defaultImplementerMaxSteps = 200

// resolveCompileOptions builds the compile options for one definition on
// one host from project config: the pinned model via the six-level
// precedence chain (models_by_host.<host> name key, then its profile key,
// then models_by_name, then models, then the compilers' preserve-merge
// fallback and host defaults), the implementer's test-file guard (default on,
// configurable), and the implementer's execution budget (agents.max_steps:
// 0 = default 200, -1 = unlimited/no field, N = N; other negatives are config
// errors).
func resolveCompileOptions(cfg *config.Config, def *subagent.Definition, hostKey string) (subagent.CompileOptions, error) {
	var opts subagent.CompileOptions
	for key := range cfg.Agents.Models {
		if !validModelProfileKeys[key] {
			return opts, fmt.Errorf("agents.models: unknown profile key %q (supported: tool-capable, tool-capable-read-only)", key)
		}
	}
	if cfg.Agents.MaxSteps < -1 {
		return opts, fmt.Errorf("agents.max_steps: invalid value %d (supported: positive budget, 0 = default %d, -1 = unlimited)", cfg.Agents.MaxSteps, defaultImplementerMaxSteps)
	}
	model, _ := resolveModelField(&cfg.Agents, def, hostKey, false)
	effort, _ := resolveModelField(&cfg.Agents, def, hostKey, true)
	opts.Model = model
	if effort != "" {
		opts.EffortConfigured = true
		if effort != "inherit" {
			opts.ReasoningEffort = effort
		}
	}
	if !cfg.Agents.DisableTestGuard && def.Name == "flowforge-implementer" {
		if len(cfg.Agents.TestFileGlobs) > 0 {
			opts.EditDeny = cfg.Agents.TestFileGlobs
		} else {
			opts.EditDeny = defaultTestFileGlobs
		}
	}
	// The question tool pause is structurally incompatible with fresh
	// execution contexts: a subagent blocking on human input freezes the
	// dispatching batch (incident: 537min session, 530min spent waiting on
	// two question calls). The only ambiguity exit is STATUS: BLOCKED.
	if def.Name == "flowforge-implementer" {
		opts.DenyQuestion = true
	}
	if def.Name == "flowforge-implementer" {
		switch cfg.Agents.MaxSteps {
		case 0:
			steps := defaultImplementerMaxSteps
			opts.MaxSteps = &steps
		case -1:
			// Explicit unlimited: no steps field, host default behavior.
		default:
			steps := cfg.Agents.MaxSteps
			opts.MaxSteps = &steps
		}
	}
	return opts, nil
}

// validateModelConfig checks the model configuration layers (agents.models,
// agents.models_by_name, agents.models_by_host) before any deploy or status
// work. Key legality is checked in full regardless of enabled hosts (a bad
// key is a corrupted-config signal); value formats are checked only against
// enabled model-carrying hosts, so a disabled host cannot block a deploy
// (design d-config-validation). Preserve-merge backfill values come from
// existing user-edited deployed files, not this command's config, and are
// not validated. It subsumes the former validateModelOverrides check
// (models_by_name keys), extended to accept profile keys alongside agent
// names. Key iteration is sorted so deploy and status report the identical
// first error.
// validateModelValueForHost checks one model value against one host's format
// rule: every host requires a non-empty token without inner whitespace or
// control characters; opencode and pi additionally require the provider/model
// shape (exactly one slash, both sides non-empty); claude accepts any single
// token (alias, inherit, or full model id). Existence is never queried
// (purely local, deterministic).
func validateModelValueForHost(where, value, host string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fmt.Errorf("%s: model value must be non-empty", where)
	}
	for _, r := range trimmed {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return fmt.Errorf("%s: model value %q must not contain whitespace or control characters", where, value)
		}
	}
	if host == "opencode" || host == "pi" {
		if strings.Count(trimmed, "/") != 1 || strings.HasPrefix(trimmed, "/") || strings.HasSuffix(trimmed, "/") {
			return fmt.Errorf("%s: model value %q invalid for host %q (want provider/model)", where, value, host)
		}
	}
	return nil
}

// validateGlobalModelValue checks a global-layer value (agents.models /
// agents.models_by_name): it compiles into every enabled model-carrying
// host, so it must satisfy each such host's format rule; a failure points
// at agents.models_by_host for host-specific configuration. The host order
// is fixed so the reported error is deterministic.
// piExtensionRelPath is the host-scoped pi extension deployed alongside
// the .pi/agents files whenever the pi host is enabled. It provides the
// implementer's test-file guard and the flowforge_frontier/flowforge_check
// native tools for the whole project, so its lifecycle is bound to the host
// selection, not to any single subagent definition.
var piExtensionRelPath = filepath.Join(".pi", "extensions", "flowforge.ts")

// deployPiExtension writes the managed flowforge extension from assets to
// .pi/extensions/flowforge.ts when pi is an enabled host. The write is
// deterministic (managed source always wins), matching the semantics of
// other host-scoped managed assets.
func deployPiExtension(projectRoot string, hosts []hostTarget) error {
	enabled := false
	for _, h := range hosts {
		if h.key == "pi" {
			enabled = true
			break
		}
	}
	if !enabled {
		return nil
	}
	assetsDir, cleanup, err := locateAssetsDir()
	if err != nil {
		return fmt.Errorf("locating assets: %w", err)
	}
	defer cleanup()
	src, err := os.ReadFile(filepath.Join(assetsDir, "pi", "flowforge.ts"))
	if err != nil {
		return fmt.Errorf("reading pi extension source: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(projectRoot, ".pi", "extensions"), 0755); err != nil {
		return fmt.Errorf("creating .pi/extensions directory: %w", err)
	}
	dst := filepath.Join(projectRoot, piExtensionRelPath)
	if err := os.WriteFile(dst, src, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", dst, err)
	}
	return nil
}

// cleanDeselectedHosts removes managed subagent files from hosts that are not
// selected. Only files matching a discoverable definition name are removed;
// project-owned files and the host directories themselves are preserved.
func cleanDeselectedHosts(projectRoot string, selected []hostTarget, allDefs []*subagent.Definition) error {
	selectedKeys := make(map[string]bool, len(selected))
	for _, h := range selected {
		selectedKeys[h.key] = true
	}
	for _, h := range allHostTargets() {
		if selectedKeys[h.key] {
			continue
		}
		dir := filepath.Join(projectRoot, h.relDir)
		for _, def := range allDefs {
			path := filepath.Join(dir, def.Name+h.ext)
			if _, err := os.Stat(path); err == nil {
				if err := os.Remove(path); err != nil {
					return fmt.Errorf("cleaning %s: %w", path, err)
				}
			}
		}
	}
	// The pi extension is a host-scoped managed resource: converge it the
	// same way per-agent files converge — present iff pi is selected.
	if !selectedKeys["pi"] {
		ext := filepath.Join(projectRoot, piExtensionRelPath)
		if _, err := os.Stat(ext); err == nil {
			if err := os.Remove(ext); err != nil {
				return fmt.Errorf("cleaning %s: %w", ext, err)
			}
		}
	}
	return nil
}

// deploySubagents discovers, compiles, and writes subagent definitions to the
// enabled host directories. If targetName is non-empty, deploys only that
// subagent. Otherwise deploys all non-disabled. Managed files in deselected
// hosts are cleaned. Before overwriting an existing deployed file, a local
// `model:` frontmatter value the fresh compile would not reproduce is
// preserved (config-pinned models always win). Returns the list of deployed
// subagent names.
func deploySubagents(projectRoot string, cfg *config.Config, targetName string) ([]string, error) {
	prepared, err := prepareSubagents(projectRoot, cfg, targetName)
	if err != nil {
		return nil, err
	}
	if len(prepared.definitions) == 0 {
		return nil, nil
	}
	for _, h := range prepared.hosts {
		if err := os.MkdirAll(filepath.Join(projectRoot, h.relDir), 0755); err != nil {
			return nil, fmt.Errorf("creating directory %s: %w", h.relDir, err)
		}
	}
	for _, output := range prepared.outputs {
		if err := os.WriteFile(output.path, output.content, 0644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", output.path, err)
		}
		for _, hint := range output.hints {
			fmt.Fprintln(os.Stderr, hint)
		}
	}
	if err := deployPiExtension(projectRoot, prepared.hosts); err != nil {
		return nil, err
	}
	if err := cleanDeselectedHosts(projectRoot, prepared.hosts, prepared.allDefs); err != nil {
		return nil, err
	}
	if err := writeAgentModelState(projectRoot, prepared.stateBytes); err != nil {
		return nil, err
	}
	var deployed []string
	for _, def := range prepared.definitions {
		deployed = append(deployed, def.Name)
	}
	return deployed, nil
}

// frontmatterModel extracts the `model` key from the leading yaml frontmatter
// block of data. Missing delimiters, unparseable yaml, or an absent/empty key
// all yield "" — preserve-merge treats such content as having no local model.
func frontmatterModel(data []byte) string {
	normalized := bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(normalized, []byte("---\n")) {
		return ""
	}
	rest := normalized[len("---\n"):]
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		return ""
	}
	var fm struct {
		Model string `yaml:"model"`
	}
	if err := yaml.Unmarshal(rest[:end], &fm); err != nil {
		return ""
	}
	return strings.TrimSpace(fm.Model)
}

// discoverSubagentSources reads subagent definitions from built-in assets and project-custom sources.
// Project-custom sources in .flowforge/subagents/ override built-in definitions with the same name.
func discoverSubagentSources(projectRoot string) ([]*subagent.Definition, error) {
	// Read built-in subagents from embedded/filesystem assets
	assetsDir, cleanup, err := locateAssetsDir()
	if err != nil {
		return nil, fmt.Errorf("locating assets: %w", err)
	}
	defer cleanup()

	builtinDir := filepath.Join(assetsDir, "subagents")
	builtinDefs, err := subagent.ParseDir(builtinDir)
	if err != nil {
		return nil, fmt.Errorf("parsing built-in subagents: %w", err)
	}

	// Read project-custom subagents from .flowforge/subagents/
	customDir := filepath.Join(projectRoot, config.ConfigDirName, "subagents")
	var customDefs []*subagent.Definition
	if stat, err := os.Stat(customDir); err == nil && stat.IsDir() {
		customDefs, err = subagent.ParseDir(customDir)
		if err != nil {
			return nil, fmt.Errorf("parsing project-custom subagents: %w", err)
		}
	}

	// Merge: custom definitions override built-in by name
	merged := make(map[string]*subagent.Definition)
	for _, def := range builtinDefs {
		merged[def.Name] = def
	}
	for _, def := range customDefs {
		merged[def.Name] = def
	}

	// Return as sorted slice
	var result []*subagent.Definition
	for _, def := range merged {
		result = append(result, def)
	}

	// Sort by name for stable output
	sortDefinitionsByName(result)
	return result, nil
}

func sortDefinitionsByName(defs []*subagent.Definition) {
	// Simple bubble sort (sufficient for small lists)
	for i := 0; i < len(defs); i++ {
		for j := i + 1; j < len(defs); j++ {
			if strings.Compare(defs[i].Name, defs[j].Name) > 0 {
				defs[i], defs[j] = defs[j], defs[i]
			}
		}
	}
}

// removeSubagent removes a subagent from all host directories.
// For built-in subagents, adds the name to config.Agents.Disabled.
// For custom subagents, deletes the source file from .flowforge/subagents/.
// Returns (isBuiltin, removedPaths, error).
func removeSubagent(projectRoot, name string) (bool, []string, error) {
	cfg, err := config.Load(projectRoot)
	if err != nil {
		return false, nil, fmt.Errorf("loading config: %w", err)
	}

	// Check if this is a built-in or custom subagent
	isBuiltin, err := isBuiltinSubagent(name)
	if err != nil {
		return false, nil, err
	}

	var removedPaths []string

	// Remove compiled files from enabled host directories only; deselected
	// hosts are converged by deploy-time cleanup.
	hosts, err := resolveHostTargets(cfg)
	if err != nil {
		return false, nil, err
	}
	for _, h := range hosts {
		path := filepath.Join(projectRoot, h.relDir, name+h.ext)
		if _, err := os.Stat(path); err == nil {
			if err := os.Remove(path); err != nil {
				return false, nil, fmt.Errorf("removing %s: %w", path, err)
			}
			removedPaths = append(removedPaths, path)
		}
	}

	if isBuiltin {
		// Add to disabled list if not already present
		alreadyDisabled := false
		for _, disabled := range cfg.Agents.Disabled {
			if disabled == name {
				alreadyDisabled = true
				break
			}
		}
		if !alreadyDisabled {
			cfg.Agents.Disabled = append(cfg.Agents.Disabled, name)
			if err := cfg.Save(projectRoot); err != nil {
				return false, nil, fmt.Errorf("saving config: %w", err)
			}
		}
	} else {
		// Custom subagent: delete source file
		sourcePath := filepath.Join(projectRoot, config.ConfigDirName, "subagents", name+".md")
		if _, err := os.Stat(sourcePath); err != nil {
			if os.IsNotExist(err) {
				return false, nil, fmt.Errorf("custom subagent source file not found: %s", sourcePath)
			}
			return false, nil, fmt.Errorf("checking source file %s: %w", sourcePath, err)
		}
		if err := os.Remove(sourcePath); err != nil {
			return false, nil, fmt.Errorf("removing source file %s: %w", sourcePath, err)
		}
		removedPaths = append(removedPaths, sourcePath)
	}

	return isBuiltin, removedPaths, nil
}

// isBuiltinSubagent checks if a subagent name exists in the built-in assets.
func isBuiltinSubagent(name string) (bool, error) {
	assetsDir, cleanup, err := locateAssetsDir()
	if err != nil {
		return false, fmt.Errorf("locating assets: %w", err)
	}
	defer cleanup()

	builtinPath := filepath.Join(assetsDir, "subagents", name+".md")
	_, err = os.Stat(builtinPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("checking built-in subagent %s: %w", name, err)
}
