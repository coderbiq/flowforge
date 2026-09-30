package command

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"strings"
)

// newModelSetCmd builds the `flowforge model-set` command group:
//
//	list          — every declared set (plus the implicit base layer),
//	                marking the currently active one
//	use <name>    — switch the active set and redeploy; on deploy failure
//	                the pointer rolls back (pointer rollback)
//	show [name]   — the effective per-agent model table for a set
//	                (default: the active set), with base/set source marks
//
// Design: docs/proposals/model-sets-switching/design.md
func newModelSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "model-set",
		Short: "Manage named agent model sets and switch between them",
	}
	cmd.AddCommand(
		newModelSetListCmd(),
		newModelSetUseCmd(),
		newModelSetShowCmd(),
	)
	return cmd
}

// loadForModelSet loads config for a model-set command, rooting it at the
// current project.
func loadForModelSet() (string, *config.Config, error) {
	projectRoot, err := config.FindProjectRoot(".")
	if err != nil {
		return "", nil, fmt.Errorf("locating project root: %w", err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		return "", nil, fmt.Errorf("loading config: %w", err)
	}
	return projectRoot, cfg, nil
}

func newModelSetListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List model sets and mark the active one",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, cfg, err := loadForModelSet()
			if err != nil {
				return err
			}
			active, err := config.ReadActiveModelSet(projectRoot)
			if err != nil {
				return fmt.Errorf("reading active model set: %w", err)
			}
			if active != "" {
				if _, ok := config.ResolveModelSet(&cfg.Agents, active); !ok {
					cmd.Printf("warning: active set %q is not declared in agents.model_sets\n\n", active)
				}
			}
			baseMark := "        "
			if active == "" {
				baseMark = "ACTIVE  "
			}
			cmd.Printf("%sdefault (base agent model layers)\n", baseMark)
			for _, name := range config.ModelSetNames(&cfg.Agents) {
				mark := "        "
				if name == active {
					mark = "ACTIVE  "
				}
				cmd.Printf("%s%s\n", mark, name)
			}
			return nil
		},
	}
}

func newModelSetUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Switch the active model set and redeploy agents",
		Long: `Switch the active model set and redeploy agents.

<name> may be a set declared in agents.model_sets, or "default" to return
to the base agent model layers. If redeploy fails, the active pointer is restored to its previous value.
Configuration, reading and compilation failures occur before artifact writes;
an I/O failure during writing can leave partially updated artifacts.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, cfg, err := loadForModelSet()
			if err != nil {
				return err
			}
			name := args[0]

			if name != "default" {
				if _, ok := config.ResolveModelSet(&cfg.Agents, name); !ok {
					return fmt.Errorf("model set %q not found in agents.model_sets (use \"default\" for the base layers)", name)
				}
			}

			prev, err := config.ReadActiveModelSet(projectRoot)
			if err != nil {
				return fmt.Errorf("reading active model set: %w", err)
			}
			if prev == name {
				cmd.Printf("model set %q is already active\n", name)
				return nil
			}

			rollback := func() {
				var rerr error
				if prev == "" {
					rerr = config.ClearActiveModelSet(projectRoot)
				} else {
					rerr = config.WriteActiveModelSet(projectRoot, prev)
				}
				if rerr != nil {
					cmd.Printf("warning: failed to restore previous model set %q: %v\n", prev, rerr)
				}
			}

			if name == "default" {
				if err := config.ClearActiveModelSet(projectRoot); err != nil {
					return fmt.Errorf("clearing active model set: %w", err)
				}
			} else if err := config.WriteActiveModelSet(projectRoot, name); err != nil {
				return fmt.Errorf("writing active model set: %w", err)
			}

			deployed, derr := deploySubagents(projectRoot, cfg, "")
			if derr != nil {
				rollback()
				return fmt.Errorf("redeploy under model set %q failed (active pointer restored to previous value): %w", name, derr)
			}
			cmd.Printf("✓ switched to model set %q; redeployed %d subagent(s)\n", name, len(deployed))
			return nil
		},
	}
}

func newModelSetShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show [name]",
		Short: "Show the effective per-agent model table for a model set",
		Long: `Show the effective per-agent model table for a model set.

[name] defaults to the active set. "default" shows the base layers. The
table applies the set overlay over the base layers; the source column
marks which layer each value comes from. Host-specific overrides
(agents.models_by_host) are listed separately below the table.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, cfg, err := loadForModelSet()
			if err != nil {
				return err
			}
			name := ""
			if len(args) == 1 {
				name = args[0]
			} else {
				if name, err = config.ReadActiveModelSet(projectRoot); err != nil {
					return fmt.Errorf("reading active model set: %w", err)
				}
			}

			var eff config.AgentsConfig
			switch {
			case name == "" || name == "default":
				name = "default"
				eff = cfg.Agents
			default:
				set, ok := config.ResolveModelSet(&cfg.Agents, name)
				if !ok {
					return fmt.Errorf("model set %q not found in agents.model_sets", name)
				}
				eff = config.ApplyModelSet(&cfg.Agents, set)
			}

			defs, err := discoverSubagentSources(projectRoot)
			if err != nil {
				return fmt.Errorf("discovering subagent sources: %w", err)
			}

			hosts, err := resolveHostTargets(cfg)
			if err != nil {
				return err
			}
			if err := validateModelConfig(cfg, defs, hosts); err != nil {
				return err
			}
			cmd.Printf("model set: %s\n\n", name)
			cmd.Printf("  HOST       AGENT                            MODEL                                    MODEL SOURCE                                         EFFORT             EFFORT SOURCE\n")
			disabled := map[string]bool{}
			for _, n := range cfg.Agents.Disabled {
				disabled[n] = true
			}
			for _, h := range hosts {
				for _, def := range defs {
					if disabled[def.Name] {
						continue
					}
					model, modelPath := resolveModelField(&eff, def, h.key, false)
					effort, effortPath := resolveModelField(&eff, def, h.key, true)
					source := func(path string, isEffort bool) string {
						if path == "" {
							return "default"
						}
						if name != "default" && declaredAt(cfg.Agents.ModelSets[name], path, isEffort) {
							return "set:" + "agents.model_sets." + name + strings.TrimPrefix(path, "agents")
						}
						return "base:" + path
					}
					modelSource, effortSource := source(modelPath, false), source(effortPath, true)
					if model == "" {
						model = "(host default)"
						if h.key == "claude" {
							model = def.ModelProfile.ClaudeModel()
						}
					}
					if effort == "inherit" {
						effort = "(host inherit)"
					} else if effort == "" {
						switch h.key {
						case "codex":
							effort = def.ModelProfile.CodexReasoningEffort()
						case "pi":
							effort = def.ModelProfile.PiThinking()
						default:
							effort = "(host default)"
						}
					}
					cmd.Printf("  %-10s %-32s %-40s %-52s %-18s %s\n", h.key, def.Name, model, modelSource, effort, effortSource)
				}
			}
			return nil
		},
	}
}

func sortedKeys(m map[string]config.ModelValue) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sortedHostKeys(m map[string]map[string]config.ModelValue) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
