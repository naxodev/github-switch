package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/naxodev/github-switch/internal/config"
	"github.com/spf13/cobra"
)

var forceRemove bool

var removeCmd = &cobra.Command{
	Use:     "remove <account-name>",
	Short:   "Remove a GitHub account",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runRemove,
}

func init() {
	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	accountName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	acc, exists := cfg.GetAccount(accountName)
	if !exists {
		return fmt.Errorf("account '%s' not found", accountName)
	}

	if !forceRemove {
		fmt.Printf("Remove account '%s'?\n", accountName)
		fmt.Printf("  Name:    %s\n", acc.Name)
		fmt.Printf("  Email:   %s\n", acc.Email)
		fmt.Printf("  SSH Key: %s\n", acc.SSHKey)
		fmt.Print("\nConfirm? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	cfg.RemoveAccount(accountName)

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Account '%s' removed.\n", accountName)
	return nil
}
