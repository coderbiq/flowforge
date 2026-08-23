package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func newCurateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "curate",
		Short: "Synthesize living documentation and archive completed proposals",
	}

	cmd.AddCommand(newCurateDiffCmd())
	cmd.AddCommand(newCurateApplyCmd())

	return cmd
}

func newCurateDiffCmd() *cobra.Command {
	var proposalID string

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Preview domain documentation diff before proposal curation",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, projectID, runtimeStore, err := currentProposalStoreWithState()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			if proposalID == "" {
				currentID, ok, _ := runtimeStore.CurrentProposalID(projectID)
				if ok && currentID != "" {
					proposalID = currentID
				}
			}

			if proposalID == "" {
				return fmt.Errorf("no active proposal specified. Use --proposal <id>")
			}

			prop, err := store.FindWorkspaceProposal(proposalID)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "=== Curate Diff Preview for Proposal: %s ===\n\n", prop.ID)

			if prop.Mode == "hierarchical" {
				modulesDir := filepath.Join(prop.Path, "modules")
				entries, err := os.ReadDir(modulesDir)
				if err == nil && len(entries) > 0 {
					fmt.Fprintln(out, "Discovered Module Updates to be synthesized into docs/domains/:")
					for _, entry := range entries {
						if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
							fmt.Fprintf(out, "  + modules/%s -> docs/domains/%s\n", entry.Name(), entry.Name())
						}
					}
					fmt.Fprintln(out)
				}
			}

			fmt.Fprintln(out, "Proposed ADR extraction target:")
			fmt.Fprintf(out, "  + docs/adrs/ADR-%s-%s.md\n\n", time.Now().Format("20060102"), prop.ID)
			fmt.Fprintln(out, "Run 'flowforge curate apply' to perform domain synthesis and archive this proposal.")
			return nil
		},
	}

	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID")
	return cmd
}

func newCurateApplyCmd() *cobra.Command {
	var proposalID string

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Synthesize domain docs, extract ADRs, and archive proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, projectID, runtimeStore, err := currentProposalStoreWithState()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			if proposalID == "" {
				currentID, ok, _ := runtimeStore.CurrentProposalID(projectID)
				if ok && currentID != "" {
					proposalID = currentID
				}
			}

			if proposalID == "" {
				return fmt.Errorf("no active proposal specified. Use --proposal <id>")
			}

			prop, err := store.FindWorkspaceProposal(proposalID)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			wikiRoot := store.WikiRootDir()

			// 1. Ensure docs/domains/ & docs/adrs/ exist
			domainsDir := filepath.Join(wikiRoot, "docs", "domains")
			adrsDir := filepath.Join(wikiRoot, "docs", "adrs")
			_ = os.MkdirAll(domainsDir, 0755)
			_ = os.MkdirAll(adrsDir, 0755)

			// 2. Synthesize hierarchical modules if present
			if prop.Mode == "hierarchical" {
				modulesDir := filepath.Join(prop.Path, "modules")
				entries, err := os.ReadDir(modulesDir)
				if err == nil {
					for _, entry := range entries {
						if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
							srcPath := filepath.Join(modulesDir, entry.Name())
							dstPath := filepath.Join(domainsDir, entry.Name())
							content, err := os.ReadFile(srcPath)
							if err == nil {
								_ = os.WriteFile(dstPath, content, 0644)
								fmt.Fprintf(out, "✓ Synced module spec to %s\n", dstPath)
							}
						}
					}
				}
			}

			// 3. Extract ADR
			adrPath := filepath.Join(adrsDir, fmt.Sprintf("ADR-%s-%s.md", time.Now().Format("20060102"), prop.ID))
			if _, err := os.Stat(adrPath); os.IsNotExist(err) {
				adrContent := fmt.Sprintf(`# ADR: %s

## Context & Problem Statement
Synthesized from proposal %s.

## Decision Drivers
- Captured consensus from workspace: %s

## Status
Accepted
`, prop.Title, prop.ID, prop.Path)
				_ = os.WriteFile(adrPath, []byte(adrContent), 0644)
				fmt.Fprintf(out, "✓ Created ADR at %s\n", adrPath)
			}

			// 4. Archive proposal (move to 02-archive/ or mark status)
			archiveDir := filepath.Join(wikiRoot, "02-archive")
			_ = os.MkdirAll(archiveDir, 0755)
			targetArchive := filepath.Join(archiveDir, prop.ID)
			if err := os.Rename(prop.Path, targetArchive); err == nil {
				fmt.Fprintf(out, "✓ Archived proposal workspace to %s\n", targetArchive)
			} else {
				// If rename fails (cross-device etc), try copying or keeping in place
				fmt.Fprintf(out, "• Completed proposal workspace kept at %s\n", prop.Path)
			}

			// Clear current proposal from runtime
			_ = clearCurrentProposalIfMatches(runtimeStore, projectID, prop.ID)

			fmt.Fprintf(out, "✓ Proposal %s successfully curated and closed.\n", prop.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID")
	return cmd
}
