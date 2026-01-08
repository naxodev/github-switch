package cmd

import (
	"fmt"

	"github.com/naxodev/github-switch/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config with default accounts",
	Long: `Initialize the configuration file with some default accounts.
This is useful for first-time setup.`,
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

	cfg.AddAccount("naxodev", config.Account{
		SSHKey: "naxo-github",
		Name:   "Nacho Vazquez",
		Email:  "nacho@naxo.dev",
	})

	cfg.AddAccount("aster", config.Account{
		SSHKey: "github_aster_secure_rsa",
		Name:   "Nacho Vazquez",
		Email:  "nacho@astercare.com",
	})

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Config initialized at: %s\n", config.GetConfigPath())
	fmt.Println("\nAccounts added:")
	for _, name := range cfg.ListAccounts() {
		acc, _ := cfg.GetAccount(name)
		fmt.Printf("  - %s (%s)\n", name, acc.Email)
	}

	return nil
}
