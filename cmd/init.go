package cmd

import (
	"fmt"

	"github.com/naxodev/github-switch/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize an empty config file",
	Long: `Initialize an empty configuration file.
Use 'github-switch add' to add your accounts after initialization.`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Accounts) > 0 {
		fmt.Println("Config already exists with accounts:")
		for _, name := range cfg.ListAccounts() {
			fmt.Printf("  - %s\n", name)
		}
		fmt.Println("\nUse 'github-switch add' to add more accounts.")
		return nil
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Config initialized at: %s\n", config.GetConfigPath())
	fmt.Println("\nUse 'github-switch add <account-name>' to add your first account.")

	return nil
}
