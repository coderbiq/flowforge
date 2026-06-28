package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/core"
	"flowforge/internal/state"
	"flowforge/internal/version"
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
  .flowforge/cache/flowforge.sqlite - Runtime state database
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
	configPath := config.ConfigPath(targetDir)

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

	if err := createConfigFile(configPath); err != nil {
		return fmt.Errorf("creating config: %w", err)
	}

	if err := createRuntimeState(targetDir); err != nil {
		return fmt.Errorf("creating runtime state: %w", err)
	}

	if err := deployManagedAssets(targetDir); err != nil {
		return fmt.Errorf("deploying managed assets: %w", err)
	}

	if err := writeProjectManifest(targetDir); err != nil {
		return fmt.Errorf("writing project manifest: %w", err)
	}

	if err := writeVersionFile(targetDir); err != nil {
		return fmt.Errorf("writing version file: %w", err)
	}

	fmt.Println("✓ FlowForge initialized successfully")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Register a project: flowforge project create <id>")
	fmt.Println("  2. Create your first proposal: flowforge proposal create <title>")
	fmt.Println("  3. Add cards: flowforge card create --type requirement --title \"...\"")

	return nil
}

func createConfigFile(configPath string) error {
	defaultConfig := config.DefaultConfig()
	content := fmt.Sprintf(`# FlowForge Configuration
version: %q

projects: []
`, defaultConfig.Version)

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

func createRuntimeState(targetDir string) error {
	cfg := config.DefaultConfig()
	dbPath := filepath.Join(cfg.CacheDir(targetDir), "flowforge.sqlite")

	store, err := state.Open(dbPath)
	if err != nil {
		return fmt.Errorf("opening runtime state store: %w", err)
	}

	if err := store.EnsureSchema(); err != nil {
		if closeErr := store.Close(); closeErr != nil {
			return fmt.Errorf("ensuring runtime state schema: %w (closing store: %v)", err, closeErr)
		}
		return fmt.Errorf("ensuring runtime state schema: %w", err)
	}

	if err := store.Close(); err != nil {
		return fmt.Errorf("closing runtime state store: %w", err)
	}

	return nil
}

func writeProjectManifest(targetDir string) error {
	manifest, err := core.GenerateManifest(embeddedAssets, version.Version)
	if err != nil {
		return fmt.Errorf("generating manifest: %w", err)
	}

	if err := manifest.Save(targetDir); err != nil {
		return fmt.Errorf("saving manifest: %w", err)
	}

	return nil
}

func writeVersionFile(targetDir string) error {
	versionPath := filepath.Join(targetDir, ".flowforge", ".version")
	if err := os.MkdirAll(filepath.Dir(versionPath), 0755); err != nil {
		return fmt.Errorf("creating version directory: %w", err)
	}
	return os.WriteFile(versionPath, []byte(version.Version+"\n"), 0644)
}
