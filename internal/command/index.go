package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

func newIndexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Manage derived sqlite indexes",
	}

	cmd.AddCommand(newIndexRebuildCmd())
	cmd.AddCommand(newIndexStatusCmd())
	cmd.AddCommand(newIndexBacklinksCmd())

	return cmd
}

func newIndexRebuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rebuild",
		Short: "Rebuild derived card indexes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, cfg, store, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(store)

			project, _, err := resolveCurrentProject(cfg, store)
			if err != nil {
				return err
			}

			wikiRoot, err := cfg.WikiRootForProject(projectRoot, project.ID)
			if err != nil {
				return err
			}

			syncSvc := state.NewCardSyncService(store.DB())
			cardStore := core.NewCardStore(wikiRoot)
			dirs := []string{
				cardStore.ActiveDir(), cardStore.IntakeDir(), cardStore.CompletedDir(),
				cardStore.LibraryDir(), cardStore.ProposalCardDir(),
			}

			indexedCards, indexedLinks, err := syncSvc.RebuildAll(cardStore.ListCardsFromFiles, dirs)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Rebuilt index for project %s\n", project.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  cards: %d\n", indexedCards)
			fmt.Fprintf(cmd.OutOrStdout(), "  links: %d\n", indexedLinks)
			return nil
		},
	}

	return cmd
}

func newIndexStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show derived index status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, cfg, store, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(store)

			project, source, err := resolveCurrentProject(cfg, store)
			if err != nil {
				return err
			}

			wikiRoot, _ := cfg.WikiRootForProject(projectRoot, project.ID)
			_ = wikiRoot

			syncSvc := state.NewCardSyncService(store.DB())

			status, err := syncSvc.DerivedIndexStatus()
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Project: %s\n", project.ID)
			fmt.Fprintf(out, "Source: %s\n", source)
			fmt.Fprintf(out, "card_index: %d\n", status.CardCount)
			fmt.Fprintf(out, "card_link: %d\n", status.LinkCount)
			return nil
		},
	}

	return cmd
}

func newIndexBacklinksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backlinks <card-id>",
		Short: "Show backlinks for a card",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]

			_, _, store, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(store)

			syncSvc := state.NewCardSyncService(store.DB())

			status, err := syncSvc.DerivedIndexStatus()
			if err != nil {
				return err
			}
			if status.CardCount == 0 && status.LinkCount == 0 {
				return fmt.Errorf("derived index is empty; run `flowforge index rebuild` first")
			}

			backlinks, err := syncSvc.Backlinks(cardID)
			if err != nil {
				return err
			}
			if len(backlinks) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No backlinks for %s\n", cardID)
				return nil
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Backlinks for %s:\n", cardID)
			for _, backlink := range backlinks {
				fmt.Fprintf(out, "%s %s\n", backlink.FromID, backlink.Relation)
			}
			return nil
		},
	}

	return cmd
}
