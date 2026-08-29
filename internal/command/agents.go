package command

import (
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

			cmd.Printf("✓ Deployed %d subagent(s) to .claude/agents/, .opencode/agent/, .codex/agents/\n", len(deployed))
			for _, name := range deployed {
				cmd.Printf("  - %s\n", name)
			}
			return nil
		},
	}

	agents.AddCommand(deploy)
	return agents
}

// deploySubagents discovers, compiles, and writes subagent definitions to host directories.
// If targetName is non-empty, deploys only that subagent. Otherwise deploys all non-disabled.
// Returns the list of deployed subagent names.
func deploySubagents(projectRoot string, cfg *config.Config, targetName string) ([]string, error) {
	// Discover sources (built-in + project-custom)
	definitions, err := discoverSubagentSources(projectRoot)
	if err != nil {
		return nil, err
	}

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

	// Create host directories
	claudeDir := filepath.Join(projectRoot, ".claude", "agents")
	opencodeDir := filepath.Join(projectRoot, ".opencode", "agent")
	codexDir := filepath.Join(projectRoot, ".codex", "agents")

	for _, dir := range []string{claudeDir, opencodeDir, codexDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	// Compile and write to each host
	var deployed []string
	for _, def := range definitions {
		// Claude Code
		claudeContent, err := subagent.CompileClaudeCode(def)
		if err != nil {
			return nil, fmt.Errorf("compiling %s for Claude Code: %w", def.Name, err)
		}
		claudePath := filepath.Join(claudeDir, def.Name+".md")
		if err := os.WriteFile(claudePath, claudeContent, 0644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", claudePath, err)
		}

		// OpenCode
		opencodeContent, err := subagent.CompileOpenCode(def)
		if err != nil {
			return nil, fmt.Errorf("compiling %s for OpenCode: %w", def.Name, err)
		}
		opencodePath := filepath.Join(opencodeDir, def.Name+".md")
		if err := os.WriteFile(opencodePath, opencodeContent, 0644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", opencodePath, err)
		}

		// Codex
		codexContent, err := subagent.CompileCodex(def)
		if err != nil {
			return nil, fmt.Errorf("compiling %s for Codex: %w", def.Name, err)
		}
		codexPath := filepath.Join(codexDir, def.Name+".toml")
		if err := os.WriteFile(codexPath, codexContent, 0644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", codexPath, err)
		}

		deployed = append(deployed, def.Name)
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
