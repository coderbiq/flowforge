package command

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/subagent"
)

// newModelSetCmd builds the `flowforge model-set` command group:
//
//	list          — every declared set (plus the implicit base layer),
//	                marking the currently active one
//	use <name>    — switch the active set and redeploy; on deploy failure
//	                the pointer rolls back (atomic switch)
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
to the base agent model layers. The switch is atomic: if the redeploy
fails, the active pointer is restored to its previous value and the deploy
error is reported.`,
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

			cmd.Printf("model set: %s\n\n", name)
			names := make([]string, 0, len(defs))
			byName := make(map[string]*subagent.Definition, len(defs))
			for _, def := range defs {
				names = append(names, def.Name)
				byName[def.Name] = def
			}
			sort.Strings(names)
			var setCfg config.ModelSetConfig
			hasSet := false
			if name != "default" {
				setCfg = cfg.Agents.ModelSets[name]
				hasSet = true
			}
			for _, n := range names {
				def := byName[n]
				model := eff.ModelOverrides[n]
				src := "base"
				if hasSet {
					if _, inSet := setCfg.ModelOverrides[n]; inSet {
						src = "set"
					}
				}
				if model == "" {
					model = eff.Models[string(def.ModelProfile)]
					src = "base"
					if hasSet {
						if _, inSet := setCfg.Models[string(def.ModelProfile)]; inSet {
							src = "set"
						}
					}
					if model == "" {
						model = "(host default)"
						src = "-"
					}
				}
				cmd.Printf("  %-32s %-40s %s\n", n, model, src)
			}
			for _, host := range sortedHostKeys(eff.ModelHostOverrides) {
				cmd.Printf("\n  models_by_host.%s:\n", host)
				inner := eff.ModelHostOverrides[host]
				for _, k := range sortedKeys(inner) {
					cmd.Printf("    %-30s %s\n", k, inner[k])
				}
			}
			return nil
		},
	}
}

func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sortedHostKeys(m map[string]map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
