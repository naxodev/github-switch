package cmd

import (
	"fmt"

	"github.com/naxodev/github-switch/internal/config"
	"github.com/naxodev/github-switch/internal/git"
	"github.com/naxodev/github-switch/internal/ssh"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current GitHub account",
	RunE:  runCurrent,
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

func runCurrent(cmd *cobra.Command, args []string) error {
	name, email, err := git.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("failed to get current Git user: %w", err)
	}

	currentKey, err := ssh.GetCurrentKey()
	if err != nil {
		return fmt.Errorf("failed to get current SSH key: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var matchedAccount string
	for accName, acc := range cfg.Accounts {
		if acc.SSHKey == currentKey || acc.Email == email {
			matchedAccount = accName
			break
		}
	}

	fmt.Println("Current configuration:")
	fmt.Printf("  Name:    %s\n", name)
	fmt.Printf("  Email:   %s\n", email)
	fmt.Printf("  SSH Key: %s\n", currentKey)

	if matchedAccount != "" {
		fmt.Printf("\nMatched account: %s\n", matchedAccount)
	} else {
		fmt.Println("\nNo matching account found in configuration.")
	}

	return nil
}
