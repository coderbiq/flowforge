package command

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

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
	cmd.AddCommand(newProposalArchiveCmd())
	cmd.AddCommand(newProposalDeleteCmd())

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

			proposalID, err := nextProposalID(store)
			if err != nil {
				return err
			}
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
	var outputFormat string
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

			if outputFormat == "json" {
				return renderProposalInspectReportJSON(cmd.OutOrStdout(), report)
			}
			return renderProposalInspectReport(cmd.OutOrStdout(), report)
		},
	}
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json")

	return cmd
}

func newProposalArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive <proposal-id>",
		Short: "Mark a proposal as completed (no directory move)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := args[0]

			store, projectID, runtimeStore, err := currentProposalStoreWithState()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			src := store.ProposalDir(proposalID)
			if err := ensureProposalDir(src, proposalID); err != nil {
				return err
			}

			propCardID := "PROP-" + proposalID
			propCard, err := store.ReadCard(propCardID)
			if err != nil {
				return fmt.Errorf("reading proposal root card: %w", err)
			}
			if propCard.Status != core.CardStatusActive {
				return fmt.Errorf("proposal %s is not active (status: %s)", proposalID, propCard.Status)
			}

			propCard.Status = core.CardStatusCompleted
			if err := store.UpdateCard(propCard); err != nil {
				return fmt.Errorf("updating proposal status: %w", err)
			}

			if err := clearCurrentProposalIfMatches(runtimeStore, projectID, proposalID); err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ Proposal %s marked as completed\n", proposalID)
			fmt.Fprintf(out, "Run 'proposal inspect %s' for final report.\n", proposalID)
			fmt.Fprintf(out, "Run 'flowforge-curate' to extract reusable knowledge to library.\n")
			return nil
		},
	}

	return cmd
}

func newProposalDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <proposal-id>",
		Short: "Delete a proposal directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := args[0]
			if !force {
				return fmt.Errorf("proposal delete requires --force")
			}

			store, projectID, runtimeStore, err := currentProposalStoreWithState()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			proposalDir, err := findProposalDirForDelete(store, proposalID)
			if err != nil {
				return err
			}

			if err := os.RemoveAll(proposalDir); err != nil {
				return fmt.Errorf("deleting proposal: %w", err)
			}
			if err := clearCurrentProposalIfMatches(runtimeStore, projectID, proposalID); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Deleted proposal %s\n", proposalID)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Delete without confirmation")

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

func currentProposalStoreWithState() (*core.CardStore, string, *state.Store, error) {
	projectRoot, cfg, runtimeStore, err := openProjectContext()
	if err != nil {
		return nil, "", nil, err
	}

	project, _, err := resolveCurrentProject(cfg, runtimeStore)
	if err != nil {
		if closeErr := runtimeStore.Close(); closeErr != nil {
			return nil, "", nil, fmt.Errorf("resolving current project: %w (closing runtime store: %v)", err, closeErr)
		}
		return nil, "", nil, err
	}

	wikiRoot, err := cfg.WikiRootForProject(projectRoot, project.ID)
	if err != nil {
		if closeErr := runtimeStore.Close(); closeErr != nil {
			return nil, "", nil, fmt.Errorf("resolving wiki root: %w (closing runtime store: %v)", err, closeErr)
		}
		return nil, "", nil, err
	}

	return core.NewCardStore(wikiRoot), project.ID, runtimeStore, nil
}

func ensureProposalDir(path string, proposalID string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("proposal %q does not exist", proposalID)
		}
		return fmt.Errorf("checking proposal dir: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("proposal %q is not a directory", proposalID)
	}

	return nil
}

func findProposalDirForDelete(store *core.CardStore, proposalID string) (string, error) {
	dir := store.FindProposalDir(proposalID)
	if dir == "" {
		return "", fmt.Errorf("proposal %q does not exist in active or completed", proposalID)
	}
	if err := ensureProposalDir(dir, proposalID); err != nil {
		return "", err
	}
	return dir, nil
}

func nextProposalID(store *core.CardStore) (string, error) {
	return nextProposalIDForPrefix(store, core.GenerateProposalIDPrefix())
}

func nextProposalIDForPrefix(store *core.CardStore, prefix string) (string, error) {
	maxSeq := 0
	for _, dir := range []string{store.ActiveDir(), store.CompletedDir()} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", fmt.Errorf("reading proposal directory %s: %w", dir, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			seq, ok := proposalSequenceForPrefix(entry.Name(), prefix)
			if ok && seq > maxSeq {
				maxSeq = seq
			}
		}
	}

	return fmt.Sprintf("%s%02d", prefix, maxSeq+1), nil
}

func proposalSequenceForPrefix(dirName string, prefix string) (int, bool) {
	if dirName == prefix {
		return 0, true
	}
	if !strings.HasPrefix(dirName, prefix) {
		return 0, false
	}

	suffix := strings.TrimPrefix(dirName, prefix)
	seqStr := suffix
	if idx := strings.Index(suffix, "_"); idx > 0 {
		seqStr = suffix[:idx]
	}
	if len(seqStr) < 2 {
		return 0, false
	}
	seq, err := strconv.Atoi(seqStr)
	if err != nil || seq < 1 {
		return 0, false
	}
	return seq, true
}

func clearCurrentProposalIfMatches(runtimeStore *state.Store, projectID string, proposalID string) error {
	currentID, ok, err := runtimeStore.CurrentProposalID(projectID)
	if err != nil {
		return err
	}
	if ok && currentID == proposalID {
		if err := runtimeStore.ClearCurrentProposalID(projectID); err != nil {
			return fmt.Errorf("clearing current proposal: %w", err)
		}
	}

	return nil
}
