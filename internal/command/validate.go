package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/core"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate cards and configuration",
		Long:  "Validate card frontmatter, configuration files, and project structure.",
	}

	cmd.AddCommand(newValidateCardCmd())
	cmd.AddCommand(newValidateAllCmd())

	return cmd
}

func newValidateCardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card <card-id-or-path>",
		Short: "Validate a card's frontmatter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]

			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			store := core.NewCardStore(cfg.WikiRoot(projectRoot))

			var filePath string
			if strings.HasSuffix(target, ".md") {
				filePath = target
				if !filepath.IsAbs(filePath) {
					filePath = filepath.Join(projectRoot, filePath)
				}
			} else {
				card, err := store.ReadCard(target)
				if err != nil {
					return fmt.Errorf("card not found: %w", err)
				}
				filePath = card.FilePath
			}

			result := core.ValidateCardFile(filePath)

			if result.HasErrors() {
				fmt.Printf("✗ Validation failed for %s:\n", filepath.Base(filePath))
				for _, e := range result.Errors {
					fmt.Printf("  - %s\n", e.Error())
				}
				os.Exit(1)
			}

			fmt.Printf("✓ Card %s is valid\n", filepath.Base(filePath))
			return nil
		},
	}

	return cmd
}

func newValidateAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Validate all cards in the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			store := core.NewCardStore(cfg.WikiRoot(projectRoot))

			dirs := []string{
				store.ActiveDir(),
				store.IntakeDir(),
				store.LibraryDir(),
			}

			totalCards := 0
			validCards := 0
			var errors []string

			for _, dir := range dirs {
				cards, err := store.ListCards(dir)
				if err != nil {
					continue
				}

				for _, card := range cards {
					totalCards++
					result := core.ValidateCardFile(card.FilePath)
					if result.HasErrors() {
						for _, e := range result.Errors {
							errors = append(errors, fmt.Sprintf("%s: %s", card.FilePath, e.Error()))
						}
					} else {
						validCards++
					}
				}
			}

			fmt.Printf("Validated %d card(s)\n", totalCards)
			fmt.Printf("  ✓ Valid: %d\n", validCards)
			fmt.Printf("  ✗ Errors: %d\n", totalCards-validCards)

			if len(errors) > 0 {
				fmt.Println("\nErrors:")
				for _, e := range errors {
					fmt.Printf("  - %s\n", e)
				}
				os.Exit(1)
			}

			return nil
		},
	}

	return cmd
}
