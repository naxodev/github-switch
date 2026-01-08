package cmd

import (
	"fmt"

	"github.com/naxodev/github-switch/internal/config"
	"github.com/naxodev/github-switch/internal/ssh"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all configured accounts",
	Aliases: []string{"ls"},
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	accounts := cfg.ListAccounts()
	if len(accounts) == 0 {
		fmt.Println("No accounts configured. Use 'github-switch add' to add an account.")
		return nil
	}

	currentKey, _ := ssh.GetCurrentKey()

	fmt.Println("Configured accounts:")
	for _, name := range accounts {
		acc, _ := cfg.GetAccount(name)
		marker := "  "
		if acc.SSHKey == currentKey {
			marker = "* "
		}
		fmt.Printf("%s%s\n", marker, name)
		fmt.Printf("    Email:   %s\n", acc.Email)
		fmt.Printf("    Name:    %s\n", acc.Name)
		fmt.Printf("    SSH Key: %s\n", acc.SSHKey)
	}

	return nil
}
