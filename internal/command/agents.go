package command

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/subagent"
)

func newAgentsCmd() *cobra.Command {
	agents := &cobra.Command{
		Use:   "agents",
		Short: "Manage subagent definitions for Claude Code, OpenCode, and Codex",
	}

	deploy := &cobra.Command{
		Use:   "deploy [name]",
		Short: "Deploy subagent definitions to host-specific directories",
		Long: `Deploy subagent definitions to .claude/agents/, .opencode/agent/, and .codex/agents/.
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
		Long: `Remove a subagent from .claude/agents/, .opencode/agent/, and .codex/agents/.
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
		Long: `Compare deployed subagent files in .claude/agents/, .opencode/agent/, and .codex/agents/
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
		{"claude", filepath.Join(".claude", "agents"), ".md", func(def *subagent.Definition, _ subagent.CompileOptions) ([]byte, error) {
			return subagent.CompileClaudeCode(def)
		}},
		{"opencode", filepath.Join(".opencode", "agent"), ".md", func(def *subagent.Definition, opts subagent.CompileOptions) ([]byte, error) {
			return subagent.CompileOpenCodeWithOptions(def, opts)
		}},
		{"codex", filepath.Join(".codex", "agents"), ".toml", func(def *subagent.Definition, _ subagent.CompileOptions) ([]byte, error) {
			return subagent.CompileCodex(def)
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
		return nil, fmt.Errorf("agents.hosts must name at least one of claude, opencode, codex")
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
			return nil, fmt.Errorf("agents.hosts: unknown host %q (supported: claude, opencode, codex)", name)
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

// resolveCompileOptions builds the OpenCode compile options for one
// definition from project config: pinned model per model profile, and the
// implementer's test-file guard (default on, configurable).
func resolveCompileOptions(cfg *config.Config, def *subagent.Definition) (subagent.CompileOptions, error) {
	var opts subagent.CompileOptions
	for key := range cfg.Agents.Models {
		if !validModelProfileKeys[key] {
			return opts, fmt.Errorf("agents.models: unknown profile key %q (supported: tool-capable, tool-capable-read-only)", key)
		}
	}
	opts.Model = cfg.Agents.Models[string(def.ModelProfile)]
	if !cfg.Agents.DisableTestGuard && def.Name == "flowforge-implementer" {
		if len(cfg.Agents.TestFileGlobs) > 0 {
			opts.EditDeny = cfg.Agents.TestFileGlobs
		} else {
			opts.EditDeny = defaultTestFileGlobs
		}
	}
	return opts, nil
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
	return nil
}

// deploySubagents discovers, compiles, and writes subagent definitions to the
// enabled host directories. If targetName is non-empty, deploys only that
// subagent. Otherwise deploys all non-disabled. Managed files in deselected
// hosts are cleaned. Returns the list of deployed subagent names.
func deploySubagents(projectRoot string, cfg *config.Config, targetName string) ([]string, error) {
	hosts, err := resolveHostTargets(cfg)
	if err != nil {
		return nil, err
	}

	// Discover sources (built-in + project-custom)
	definitions, err := discoverSubagentSources(projectRoot)
	if err != nil {
		return nil, err
	}
	allDefs := definitions

	// Filter by targetName if specified
	if targetName != "" {
		found := false
		for _, def := range definitions {
			if def.Name == targetName {
				definitions = []*subagent.Definition{def}
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("subagent %q not found in built-in or project sources", targetName)
		}
	} else {
		// Filter out disabled subagents
		disabled := make(map[string]bool)
		for _, name := range cfg.Agents.Disabled {
			disabled[name] = true
		}
		var filtered []*subagent.Definition
		for _, def := range definitions {
			if !disabled[def.Name] {
				filtered = append(filtered, def)
			}
		}
		definitions = filtered
	}

	if len(definitions) == 0 {
		return nil, nil
	}

	// Create enabled host directories
	for _, h := range hosts {
		if err := os.MkdirAll(filepath.Join(projectRoot, h.relDir), 0755); err != nil {
			return nil, fmt.Errorf("creating directory %s: %w", h.relDir, err)
		}
	}

	// Compile and write to each enabled host
	var deployed []string
	for _, def := range definitions {
		opts, err := resolveCompileOptions(cfg, def)
		if err != nil {
			return nil, err
		}
		for _, h := range hosts {
			content, err := h.compile(def, opts)
			if err != nil {
				return nil, fmt.Errorf("compiling %s for %s: %w", def.Name, h.key, err)
			}
			path := filepath.Join(projectRoot, h.relDir, def.Name+h.ext)
			if err := os.WriteFile(path, content, 0644); err != nil {
				return nil, fmt.Errorf("writing %s: %w", path, err)
			}
		}
		deployed = append(deployed, def.Name)
	}

	if err := cleanDeselectedHosts(projectRoot, hosts, allDefs); err != nil {
		return nil, err
	}

	return deployed, nil
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
