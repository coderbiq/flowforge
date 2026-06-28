package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"flowforge/internal/state"
	"flowforge/internal/update"
	"flowforge/internal/version"
)

var cfgFile string
var noVersionCheck bool

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flowforge",
		Short: "FlowForge — AI-assisted software design & delivery toolkit",
		Long: `FlowForge is a workflow toolkit for AI-assisted software design and delivery.
It provides card-based knowledge management, task orchestration, and context aggregation
through a CLI-first interface.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default: $HOME/.flowforge/config.yaml)")
	cmd.PersistentFlags().StringP("output", "o", "text",
		"output format: text, json")
	cmd.PersistentFlags().BoolVar(&noVersionCheck, "no-version-check", false,
		"skip automatic version check")

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newProjectCmd())
	cmd.AddCommand(newCardCmd())
	cmd.AddCommand(newTaskCmd())
	cmd.AddCommand(newProposalCmd())
	cmd.AddCommand(newIndexCmd())
	cmd.AddCommand(newLibraryCmd())
	cmd.AddCommand(newContextCmd())
	cmd.AddCommand(newLogCmd())
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newStructureCmd())
	cmd.AddCommand(newSkillCmd())
	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newUpgradeCmd())
	cmd.AddCommand(newUninstallCmd())

	return cmd
}

func initConfig(cmd *cobra.Command) error {
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("finding home dir: %w", err)
		}
		viper.AddConfigPath(home + "/.flowforge")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("FLOWFORGE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}

	_ = viper.ReadInConfig()

	if noVersionCheck || cmd.Name() == "version" {
		return nil
	}

	if !viper.GetBool("version_check") && viper.IsSet("version_check") {
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	dbPath := filepath.Join(home, ".flowforge", "cache", "flowforge.sqlite")
	store, err := state.Open(dbPath)
	if err != nil {
		return nil
	}

	if err := store.EnsureSchema(); err != nil {
		store.Close()
		return nil
	}

	checker := update.NewVersionChecker(version.Version, store)
	checker.CheckAsync(func(msg string) {
		fmt.Fprintf(cmd.ErrOrStderr(), "\n%s\n", msg)
	})

	return nil
}
