package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/naxodev/github-switch/internal/config"
	"github.com/naxodev/github-switch/internal/git"
	"github.com/naxodev/github-switch/internal/ssh"
	"github.com/spf13/cobra"
)

var (
	forceSwitch bool
)

var switchCmd = &cobra.Command{
	Use:   "switch [account]",
	Short: "Switch to a GitHub account",
	Long: `Switch to a different GitHub account by updating SSH config
and global Git configuration.

If no account is specified, an interactive menu will be shown.`,
	Aliases: []string{"sw"},
	RunE:    runSwitch,
}

func init() {
	switchCmd.Flags().BoolVarP(&forceSwitch, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(switchCmd)
}

func runSwitch(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Accounts) == 0 {
		return fmt.Errorf("no accounts configured. Use 'github-switch add' to add an account")
	}

	var accountName string
	if len(args) > 0 {
		accountName = args[0]
	} else {
		accountName, err = selectAccount(cfg)
		if err != nil {
			return err
		}
	}

	account, ok := cfg.GetAccount(accountName)
	if !ok {
		return fmt.Errorf("unknown account: %s", accountName)
	}

	if !forceSwitch {
		fmt.Printf("Switch to account '%s'?\n", accountName)
		fmt.Printf("  Name:    %s\n", account.Name)
		fmt.Printf("  Email:   %s\n", account.Email)
		fmt.Printf("  SSH Key: %s\n", account.SSHKey)
		fmt.Print("\nConfirm? [Y/n]: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "" && response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	if err := ssh.UpdateConfig(account.SSHKey); err != nil {
		return fmt.Errorf("failed to update SSH config: %w", err)
	}

	if err := git.UpdateGlobalConfig(account.Name, account.Email); err != nil {
		return fmt.Errorf("failed to update Git config: %w", err)
	}

	if err := ssh.AddKeyToAgent(account.SSHKey); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to add key to ssh-agent: %v\n", err)
	}

	fmt.Printf("Switched to GitHub account: %s\n", accountName)
	return nil
}

func selectAccount(cfg *config.Config) (string, error) {
	accounts := cfg.ListAccounts()
	if len(accounts) == 0 {
		return "", fmt.Errorf("no accounts available")
	}

	fmt.Println("Select an account:")
	for i, name := range accounts {
		acc, _ := cfg.GetAccount(name)
		fmt.Printf("  %d. %s (%s)\n", i+1, name, acc.Email)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter your choice (1-%d): ", len(accounts))
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		selection, err := strconv.Atoi(input)
		if err == nil && selection > 0 && selection <= len(accounts) {
			return accounts[selection-1], nil
		}

		for _, name := range accounts {
			if strings.EqualFold(input, name) {
				return name, nil
			}
		}

		fmt.Println("Invalid selection. Please try again.")
	}
}
