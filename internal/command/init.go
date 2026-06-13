package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
)

func newInitCmd() *cobra.Command {
	var (
		yes      bool
		template string
	)

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize FlowForge in the current or specified directory",
		Long: `Initialize FlowForge project structure in the target directory.

This creates:
  .flowforge/config.yaml    - Project configuration
  ff-wiki/                  - Wiki root directory
    ├── 01-workspace/       - Active proposals
    │   ├── 01-active/       - Current proposals
    │   ├── 02-intake/       - Pending requirements
    │   └── 03-completed/    - Archived proposals
    └── 02-library/         - Archived knowledge
        ├── 10-requirements/
        ├── 20-decisions/
        ├── 30-designs/
        ├── 40-tasks/
        ├── 50-logs/
        ├── 60-conventions/
        ├── 70-findings/
        └── 80-modules/
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir := "."
			if len(args) > 0 {
				targetDir = args[0]
			}

			absDir, err := filepath.Abs(targetDir)
			if err != nil {
				return fmt.Errorf("resolving path: %w", err)
			}

			return runInit(absDir, yes, template)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompts")
	cmd.Flags().StringVar(&template, "template", "default", "Project template (default/minimal)")

	return cmd
}

func runInit(targetDir string, yes bool, template string) error {
	configDir := filepath.Join(targetDir, config.ConfigDirName)
	configPath := filepath.Join(configDir, config.ConfigFileName)

	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("FlowForge already initialized in %s (config.yaml exists)", targetDir)
	}

	if !yes {
		fmt.Printf("Initialize FlowForge in %s? [y/N] ", targetDir)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	if err := createDirectoryStructure(targetDir); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	if err := createConfigFile(configPath); err != nil {
		return fmt.Errorf("creating config: %w", err)
	}

	if err := createHomeIndex(targetDir); err != nil {
		return fmt.Errorf("creating home index: %w", err)
	}

	fmt.Println("✓ FlowForge initialized successfully")
	fmt.Printf("  Config: %s\n", configPath)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Edit .flowforge/config.yaml to customize project settings")
	fmt.Println("  2. Create your first proposal: flowforge proposal create <title>")
	fmt.Println("  3. Add cards: flowforge card create --type requirement --title \"...\"")

	return nil
}

func createDirectoryStructure(targetDir string) error {
	dirs := []string{
		filepath.Join(targetDir, config.ConfigDirName),
		filepath.Join(targetDir, config.ConfigDirName, "cache"),
		filepath.Join(targetDir, "ff-wiki"),
		filepath.Join(targetDir, "ff-wiki", "01-workspace"),
		filepath.Join(targetDir, "ff-wiki", "01-workspace", "01-active"),
		filepath.Join(targetDir, "ff-wiki", "01-workspace", "02-intake"),
		filepath.Join(targetDir, "ff-wiki", "01-workspace", "03-completed"),
		filepath.Join(targetDir, "ff-wiki", "02-library"),
		filepath.Join(targetDir, "ff-wiki", "02-library", "10-requirements"),
		filepath.Join(targetDir, "ff-wiki", "02-library", "20-decisions"),
		filepath.Join(targetDir, "ff-wiki", "02-library", "30-designs"),
		filepath.Join(targetDir, "ff-wiki", "02-library", "40-tasks"),
		filepath.Join(targetDir, "ff-wiki", "02-library", "50-logs"),
		filepath.Join(targetDir, "ff-wiki", "02-library", "60-conventions"),
		filepath.Join(targetDir, "ff-wiki", "02-library", "70-findings"),
		filepath.Join(targetDir, "ff-wiki", "02-library", "80-modules"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", dir, err)
		}
	}

	return nil
}

func createConfigFile(configPath string) error {
	content := `# FlowForge Configuration
version: "2.0.0"

projects:
  - id: "default"
    wikiRoot: "ff-wiki"
    srcDirs: []
`

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

func createHomeIndex(targetDir string) error {
	homeIndexPath := filepath.Join(targetDir, "ff-wiki", "00-STR-HOME.md")

	content := `---
id: STR-HOME
title: "FlowForge Knowledge Base"
type: structure
status: active
cards: []
---

# FlowForge Knowledge Base

Welcome to your FlowForge knowledge base.

## Structure

- **01-workspace/01-active/** - Current proposals and their cards
- **01-workspace/02-intake/** - Pending requirements awaiting triage
- **01-workspace/03-completed/** - Archived proposals
- **02-library/** - Archived knowledge organized by type

## Getting Started

1. Create a proposal: ` + "`flowforge proposal create \"My Feature\"`" + `
2. Add cards to the proposal: ` + "`flowforge card create --type requirement --title \"...\"`" + `
3. Track progress: ` + "`flowforge card list --status in_progress`" + `

## Card Types

- **requirement** - User needs and features
- **decision** - Architecture and design decisions
- **design** - Technical designs and specifications
- **task** - Implementation tasks
- **log** - Work logs and progress notes
- **convention** - Coding standards and conventions
- **finding** - Discoveries and insights
- **module** - Module documentation
- **structure** - Index cards organizing related content
`

	if err := os.WriteFile(homeIndexPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing home index: %w", err)
	}

	return nil
}
