package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newMemoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memory",
		Short: "Manage Tier 1 (global) and Tier 2 (proposal) working memory",
	}

	cmd.AddCommand(newMemoryInitCmd())
	cmd.AddCommand(newMemoryShowCmd())

	return cmd
}

func newMemoryInitCmd() *cobra.Command {
	var proposalID string
	var hierarchical bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize global memory (docs/CONTEXT.md) or proposal workspace (01-workspace/<proposal_id>)",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, projectID, runtimeStore, err := currentProposalStoreWithState()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			out := cmd.OutOrStdout()

			if proposalID == "" {
				// Initialize Tier 1 Global Memory: docs/CONTEXT.md
				wikiRoot := store.WikiRootDir()
				docsDir := filepath.Join(wikiRoot, "docs")
				if err := os.MkdirAll(docsDir, 0755); err != nil {
					return fmt.Errorf("creating docs directory: %w", err)
				}
				contextFile := filepath.Join(docsDir, "CONTEXT.md")
				if _, err := os.Stat(contextFile); os.IsNotExist(err) {
					template := `# Project Context (Tier 1 Working Memory)

## Ubiquitous Language
- **Domain**: Key terms and definitions.

## Architectural Boundaries
- Enforced module boundaries and forbidden dependencies.

## Active Proposals
- (Run 'flowforge status' to view ongoing work)
`
					if err := os.WriteFile(contextFile, []byte(template), 0644); err != nil {
						return fmt.Errorf("writing global context file: %w", err)
					}
					fmt.Fprintf(out, "✓ Initialized Tier 1 Global Memory at %s\n", contextFile)
				} else {
					fmt.Fprintf(out, "• Tier 1 Global Memory already exists at %s\n", contextFile)
				}
				return nil
			}

			// Initialize Tier 2 Proposal Working Memory
			workspaceDir := store.WorkspaceDir()
			propDir := filepath.Join(workspaceDir, proposalID)
			if err := os.MkdirAll(propDir, 0755); err != nil {
				return fmt.Errorf("creating proposal directory: %w", err)
			}

			if hierarchical {
				modulesDir := filepath.Join(propDir, "modules")
				refsDir := filepath.Join(propDir, "references")
				if err := os.MkdirAll(modulesDir, 0755); err != nil {
					return fmt.Errorf("creating modules directory: %w", err)
				}
				if err := os.MkdirAll(refsDir, 0755); err != nil {
					return fmt.Errorf("creating references directory: %w", err)
				}
			}

			readmePath := filepath.Join(propDir, "README.md")
			if _, err := os.Stat(readmePath); os.IsNotExist(err) {
				var template strings.Builder
				template.WriteString(fmt.Sprintf("# Proposal: %s\n\n", proposalID))
				template.WriteString("Status: active\n\n")
				template.WriteString("## Objective\nBrief description of proposal objective.\n\n")
				template.WriteString("## Grilling & Consensus\n- [ ] Key boundaries clarified\n\n")
				template.WriteString("## Plan & Slices\n")
				template.WriteString("### Slice 1: Tracer Bullet\n")
				template.WriteString("- **Objective**: End-to-end walking skeleton\n")
				template.WriteString("- **Touchpoints**: `src/`\n")
				template.WriteString("- **Test**: `go test ./...`\n\n")

				if err := os.WriteFile(readmePath, []byte(template.String()), 0644); err != nil {
					return fmt.Errorf("writing proposal README: %w", err)
				}
				fmt.Fprintf(out, "✓ Initialized Tier 2 Working Memory for proposal %s at %s\n", proposalID, readmePath)
			} else {
				fmt.Fprintf(out, "• Proposal working memory already exists at %s\n", readmePath)
			}

			// Record active proposal in runtime store
			_ = runtimeStore.SetCurrentProposalID(projectID, proposalID)
			return nil
		},
	}

	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (e.g. CR26082201_feature)")
	cmd.Flags().BoolVar(&hierarchical, "hierarchical", false, "Create hierarchical module structure (modules/, references/)")

	return cmd
}

func newMemoryShowCmd() *cobra.Command {
	var proposalID string

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Display working memory summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, projectID, runtimeStore, err := currentProposalStoreWithState()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			out := cmd.OutOrStdout()

			if proposalID == "" {
				currentID, ok, _ := runtimeStore.CurrentProposalID(projectID)
				if ok && currentID != "" {
					proposalID = currentID
				}
			}

			if proposalID == "" {
				// Show Tier 1 Global Memory
				contextFile := filepath.Join(store.WikiRootDir(), "docs", "CONTEXT.md")
				content, err := os.ReadFile(contextFile)
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Fprintln(out, "No global memory (docs/CONTEXT.md) found. Run 'flowforge memory init' to create one.")
						return nil
					}
					return err
				}
				fmt.Fprintf(out, "=== Tier 1: Global Working Memory (%s) ===\n\n%s\n", contextFile, string(content))
				return nil
			}

			// Show Tier 2 Proposal Working Memory
			prop, err := store.FindWorkspaceProposal(proposalID)
			if err != nil {
				return err
			}

			content, err := os.ReadFile(prop.Readme)
			if err != nil {
				return fmt.Errorf("reading proposal readme: %w", err)
			}

			fmt.Fprintf(out, "=== Tier 2: Proposal Scratchpad [%s] (%s) ===\nMode: %s | Status: %s\n\n%s\n",
				prop.ID, prop.Path, prop.Mode, prop.Status, string(content))
			return nil
		},
	}

	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID")

	return cmd
}
