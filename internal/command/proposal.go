package command

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

func newProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal",
		Short: "Manage proposals",
	}

	cmd.AddCommand(newProposalCreateCmd())
	cmd.AddCommand(newProposalUseCmd())
	cmd.AddCommand(newProposalCurrentCmd())
	cmd.AddCommand(newProposalListCmd())
	cmd.AddCommand(newProposalInspectCmd())

	return cmd
}

func newProposalCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <title>",
		Short: "Create a new proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]

			projectRoot, cfg, runtimeStore, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			project, _, err := resolveCurrentProject(cfg, runtimeStore)
			if err != nil {
				return err
			}

			wikiRoot, err := cfg.WikiRootForProject(projectRoot, project.ID)
			if err != nil {
				return err
			}
			store := core.NewCardStore(wikiRoot)

			proposalID := core.GenerateProposalID()
			rootPath, indexPath, err := store.CreateProposal(proposalID, title)
			if err != nil {
				return err
			}
			if err := runtimeStore.SetCurrentProposalID(project.ID, proposalID); err != nil {
				return fmt.Errorf("setting current proposal: %w", err)
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ Created proposal %s\n", proposalID)
			fmt.Fprintf(out, "  Title: %s\n", title)
			fmt.Fprintf(out, "  Root card: %s\n", rootPath)
			fmt.Fprintf(out, "  Requirement index: %s\n", indexPath)
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Next steps:")
			fmt.Fprintf(out, "  flowforge card create --type requirement --title \"...\" --proposal %s\n", proposalID)

			return nil
		},
	}

	return cmd
}

func newProposalUseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use <proposal-id>",
		Short: "Set the current proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := args[0]

			projectRoot, cfg, runtimeStore, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			project, _, err := resolveCurrentProject(cfg, runtimeStore)
			if err != nil {
				return err
			}

			wikiRoot, err := cfg.WikiRootForProject(projectRoot, project.ID)
			if err != nil {
				return err
			}
			store := core.NewCardStore(wikiRoot)

			proposalDir := store.ProposalDir(proposalID)
			info, err := os.Stat(proposalDir)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("proposal %q does not exist in current project %q", proposalID, project.ID)
				}
				return fmt.Errorf("checking proposal dir: %w", err)
			}
			if !info.IsDir() {
				return fmt.Errorf("proposal %q is not a directory", proposalID)
			}

			if err := runtimeStore.SetCurrentProposalID(project.ID, proposalID); err != nil {
				return fmt.Errorf("setting current proposal: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Current proposal: %s\n", proposalID)
			return nil
		},
	}

	return cmd
}

func newProposalCurrentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show the current proposal",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, cfg, runtimeStore, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			project, _, err := resolveCurrentProject(cfg, runtimeStore)
			if err != nil {
				return err
			}

			proposalID, err := currentProposalIDForProject(runtimeStore, project.ID)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Proposal: %s\n", proposalID)
			fmt.Fprintf(cmd.OutOrStdout(), "Project: %s\n", project.ID)
			return nil
		},
	}

	return cmd
}

func newProposalListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active proposals",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, cfg, runtimeStore, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			project, _, err := resolveCurrentProject(cfg, runtimeStore)
			if err != nil {
				return err
			}

			wikiRoot, err := cfg.WikiRootForProject(projectRoot, project.ID)
			if err != nil {
				return err
			}
			store := core.NewCardStore(wikiRoot)

			entries, err := os.ReadDir(store.ActiveDir())
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintln(cmd.OutOrStdout(), "No active proposals.")
					return nil
				}
				return fmt.Errorf("reading active proposals: %w", err)
			}

			var proposalIDs []string
			for _, entry := range entries {
				if entry.IsDir() {
					proposalIDs = append(proposalIDs, entry.Name())
				}
			}
			sort.Strings(proposalIDs)

			out := cmd.OutOrStdout()
			if len(proposalIDs) == 0 {
				fmt.Fprintln(out, "No active proposals.")
				return nil
			}

			fmt.Fprintln(out, "Active proposals:")
			for _, proposalID := range proposalIDs {
				fmt.Fprintf(out, "- %s\n", proposalID)
			}
			return nil
		},
	}

	return cmd
}

func newProposalInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect <proposal-id>",
		Short: "Inspect a proposal summary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := currentCardStore()
			if err != nil {
				return err
			}

			report, err := buildProposalInspectReport(store, args[0])
			if err != nil {
				return err
			}

			return renderProposalInspectReport(cmd.OutOrStdout(), report)
		},
	}

	return cmd
}

func currentProposalIDForProject(runtimeStore *state.Store, projectID string) (string, error) {
	if runtimeStore == nil {
		return "", fmt.Errorf("runtime state store is required")
	}
	proposalID, ok, err := runtimeStore.CurrentProposalID(projectID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("no current proposal set for project %q; run flowforge proposal use <proposal-id>", projectID)
	}

	return proposalID, nil
}
