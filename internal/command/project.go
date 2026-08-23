package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/core"
	"flowforge/internal/state"
)

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage FlowForge projects",
	}

	cmd.AddCommand(newProjectCreateCmd())
	cmd.AddCommand(newProjectListCmd())
	cmd.AddCommand(newProjectCurrentCmd())
	cmd.AddCommand(newProjectUseCmd())

	return cmd
}

func newProjectCreateCmd() *cobra.Command {
	var (
		wikiRoot     string
		srcDirs      []string
		setAsDefault bool
	)

	cmd := &cobra.Command{
		Use:   "create <project-id>",
		Short: "Register a project and bootstrap its wiki root",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectID := args[0]
			projectRoot, cfg, store, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(store)

			if _, ok := cfg.ProjectByID(projectID); ok {
				return fmt.Errorf("project %q already exists", projectID)
			}

			resolvedWikiRoot := wikiRoot
			if resolvedWikiRoot == "" {
				resolvedWikiRoot = defaultWikiRootForProject(projectID)
			}

			project := config.ProjectConfig{
				ID:       projectID,
				WikiRoot: resolvedWikiRoot,
				SrcDirs:  append([]string(nil), srcDirs...),
			}

			absoluteWikiRoot := resolveProjectWikiRoot(projectRoot, project)
			if err := createProjectWikiRoot(absoluteWikiRoot, projectID, project.SrcDirs); err != nil {
				return err
			}

			cfg.Projects = append(cfg.Projects, project)
			if err := cfg.Save(projectRoot); err != nil {
				return err
			}

			if setAsDefault || len(cfg.Projects) == 1 {
				if err := store.SetCurrentProjectID(projectID); err != nil {
					return fmt.Errorf("setting current project: %w", err)
				}
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ Project created: %s\n", projectID)
			if len(project.SrcDirs) > 0 {
				fmt.Fprintf(out, "  srcDirs: %v\n", project.SrcDirs)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&wikiRoot, "wiki-root", "", "Wiki root directory for the project")
	cmd.Flags().StringSliceVar(&srcDirs, "src-dir", nil, "Project source directory (repeatable)")
	cmd.Flags().BoolVar(&setAsDefault, "default", false, "Set the new project as current")

	return cmd
}

func newProjectListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List registered projects",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, cfg, store, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(store)

			currentID, _, err := store.CurrentProjectID()
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if len(cfg.Projects) == 0 {
				fmt.Fprintln(out, "No projects registered.")
				fmt.Fprintln(out, "Run: flowforge project create <id>")
				return nil
			}

			fmt.Fprintln(out, "Projects:")
			for _, project := range cfg.Projects {
				marker := " "
				if project.ID == currentID {
					marker = "*"
				}

				_, err := cfg.WikiRootForProject(projectRoot, project.ID)
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "%s %s\n", marker, project.ID)
				if len(project.SrcDirs) > 0 {
					fmt.Fprintf(out, "    srcDirs: %v\n", project.SrcDirs)
				}
			}

			return nil
		},
	}

	return cmd
}

func newProjectCurrentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show the current project",
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

			_, err = cfg.WikiRootForProject(projectRoot, project.ID)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Project: %s\n", project.ID)
			fmt.Fprintf(out, "Source: %s\n", source)
			if len(project.SrcDirs) > 0 {
				fmt.Fprintf(out, "SrcDirs: %v\n", project.SrcDirs)
			}

			return nil
		},
	}

	return cmd
}

func newProjectUseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use <project-id>",
		Short: "Set the current project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectID := args[0]
			_, cfg, store, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(store)

			if _, ok := cfg.ProjectByID(projectID); !ok {
				return fmt.Errorf("project %q is not registered", projectID)
			}

			if err := store.SetCurrentProjectID(projectID); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Current project: %s\n", projectID)
			return nil
		},
	}

	return cmd
}

func openProjectContext() (string, *config.Config, *state.Store, error) {
	svc, err := currentConfigService()
	if err != nil {
		return "", nil, nil, err
	}
	return svc.ProjectRoot(), svc.FileStore().Config(), svc.StateStore().DB(), nil
}

func currentConfigService() (*config.ConfigService, error) {
	return config.New(".")
}

func closeStateStore(store *state.Store) {
	if store == nil {
		return
	}
	if err := store.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: closing runtime state store: %v\n", err)
	}
}

func resolveCurrentProject(cfg *config.Config, store *state.Store) (config.ProjectConfig, string, error) {
	if cfg == nil {
		return config.ProjectConfig{}, "", fmt.Errorf("config is required")
	}
	if store == nil {
		return config.ProjectConfig{}, "", fmt.Errorf("runtime state store is required")
	}

	projectID, ok, err := store.CurrentProjectID()
	if err != nil {
		return config.ProjectConfig{}, "", err
	}
	if ok {
		project, found := cfg.ProjectByID(projectID)
		if !found {
			return config.ProjectConfig{}, "", fmt.Errorf("current project %q is not registered", projectID)
		}
		return project, "runtime-state", nil
	}

	if len(cfg.Projects) == 1 {
		return cfg.Projects[0], "single-project", nil
	}

	if len(cfg.Projects) == 0 {
		return config.ProjectConfig{}, "", fmt.Errorf("no projects registered; run flowforge project create <id>")
	}

	return config.ProjectConfig{}, "", fmt.Errorf("current project is not set; run flowforge project use <id>")
}

func currentCardStore() (*core.CardStore, error) {
	svc, err := currentConfigService()
	if err != nil {
		return nil, err
	}

	project, _, err := resolveCurrentProject(svc.FileStore().Config(), svc.StateStore().DB())
	if err != nil {
		svc.Close()
		return nil, err
	}

	wikiRoot, err := svc.WikiRoot(project.ID)
	if err != nil {
		svc.Close()
		return nil, err
	}

	syncSvc := state.NewCardSyncService(svc.StateStore().DB().DB())
	cardStore := core.NewCardStoreWithSync(wikiRoot, syncSvc)
	return cardStore, nil
}

func resolveDefaultProposalID(explicitProposalID string, cardType core.CardType) (string, error) {
	if explicitProposalID != "" || !cardTypeDefaultsToCurrentProposal(cardType) {
		return explicitProposalID, nil
	}

	svc, err := currentConfigService()
	if err != nil {
		return "", err
	}
	defer svc.Close()

	project, _, err := resolveCurrentProject(svc.FileStore().Config(), svc.StateStore().DB())
	if err != nil {
		return "", err
	}

	proposalID, ok, err := svc.CurrentProposalID(project.ID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}

	return proposalID, nil
}

func cardTypeDefaultsToCurrentProposal(cardType core.CardType) bool {
	switch cardType {
	case core.CardTypeDecision, core.CardTypeFinding:
		return true
	default:
		return false
	}
}

func defaultWikiRootForProject(projectID string) string {
	if projectID == "default" {
		return "ff-wiki"
	}

	return fmt.Sprintf("ff-wiki-%s", projectID)
}

func resolveProjectWikiRoot(projectRoot string, project config.ProjectConfig) string {
	if filepath.IsAbs(project.WikiRoot) {
		return project.WikiRoot
	}

	return filepath.Join(projectRoot, project.WikiRoot)
}

func createProjectWikiRoot(wikiRoot string, projectID string, srcDirs []string) error {
	dirs := []string{
		wikiRoot,
		filepath.Join(wikiRoot, "01-workspace"),
		filepath.Join(wikiRoot, "02-archive"),
		filepath.Join(wikiRoot, "02-library"),
		filepath.Join(wikiRoot, "02-library", "20-decisions"),
		filepath.Join(wikiRoot, "02-library", "60-conventions"),
		filepath.Join(wikiRoot, "02-library", "70-findings"),
		filepath.Join(wikiRoot, "02-library", "80-modules"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", dir, err)
		}
	}

	if err := writeProjectHomeIndex(wikiRoot, projectID, srcDirs); err != nil {
		return err
	}

	return nil
}

func writeProjectHomeIndex(wikiRoot string, projectID string, srcDirs []string) error {
	homeIndexPath := filepath.Join(wikiRoot, "00-FLOWFORGE-HOME.md")
	proposalCmd := "`flowforge memory init --proposal <id>`"
	sliceCmd := "`flowforge context slice --slice <n>`"
	statusCmd := "`flowforge status`"

	content := fmt.Sprintf(`---
title: "FlowForge Knowledge Base"
---

# FlowForge Knowledge Base

Project: %s

## Structure

- **01-workspace/** - Active proposal scratchpads & working memory
- **02-archive/** - Completed & curated proposals
- **02-library/** - Archived knowledge organized by type
- **docs/domains/** - Living domain specifications

## Getting Started

1. Initialize proposal working memory: %s
2. Extract slice execution context: %s
3. Track overall progress: %s
`, projectID, proposalCmd, sliceCmd, statusCmd)

	if len(srcDirs) > 0 {
		content += "\n## Source Directories\n\n"
		for _, srcDir := range srcDirs {
			content += fmt.Sprintf("- %s\n", srcDir)
		}
	}

	if err := os.WriteFile(homeIndexPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing home index: %w", err)
	}

	return nil
}
