package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/naxodev/github-switch/internal/config"
	"github.com/spf13/cobra"
)

var (
	addName   string
	addEmail  string
	addSSHKey string
)

var addCmd = &cobra.Command{
	Use:   "add <account-name>",
	Short: "Add a new GitHub account",
	Long: `Add a new GitHub account configuration.

You can specify options via flags or interactively.`,
	Args: cobra.ExactArgs(1),
	RunE: runAdd,
}

func init() {
	addCmd.Flags().StringVarP(&addName, "name", "n", "", "Git user name")
	addCmd.Flags().StringVarP(&addEmail, "email", "e", "", "Git email address")
	addCmd.Flags().StringVarP(&addSSHKey, "ssh-key", "k", "", "SSH key filename (in ~/.ssh/)")
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	accountName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if _, exists := cfg.GetAccount(accountName); exists {
		return fmt.Errorf("account '%s' already exists", accountName)
	}

	reader := bufio.NewReader(os.Stdin)

	if addName == "" {
		fmt.Print("Git user name: ")
		addName, _ = reader.ReadString('\n')
		addName = strings.TrimSpace(addName)
	}

	if addEmail == "" {
		fmt.Print("Git email: ")
		addEmail, _ = reader.ReadString('\n')
		addEmail = strings.TrimSpace(addEmail)
	}

	if addSSHKey == "" {
		availableKeys, _ := listSSHKeys()
		if len(availableKeys) > 0 {
			fmt.Println("Available SSH keys:")
			for i, key := range availableKeys {
				fmt.Printf("  %d. %s\n", i+1, key)
			}
		}
		fmt.Print("SSH key filename: ")
		addSSHKey, _ = reader.ReadString('\n')
		addSSHKey = strings.TrimSpace(addSSHKey)
	}

	if addName == "" || addEmail == "" || addSSHKey == "" {
		return fmt.Errorf("all fields are required")
	}

	cfg.AddAccount(accountName, config.Account{
		Name:   addName,
		Email:  addEmail,
		SSHKey: addSSHKey,
	})

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Account '%s' added successfully.\n", accountName)
	fmt.Printf("Config saved to: %s\n", config.GetConfigPath())
	return nil
}

func listSSHKeys() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sshDir := filepath.Join(home, ".ssh")
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".pub") || name == "config" || name == "known_hosts" || name == "authorized_keys" {
			continue
		}
		if strings.Contains(name, "rsa") || strings.Contains(name, "ed25519") || strings.Contains(name, "ecdsa") || strings.HasPrefix(name, "id_") || strings.Contains(name, "github") {
			keys = append(keys, name)
		}
	}

	return keys, nil
}
