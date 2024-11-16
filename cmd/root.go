package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "github-switch [account]",
	Short: "Switch between GitHub accounts",
	Long: `A simple CLI tool to switch between different GitHub accounts
by modifying SSH config and Git configuration.`,
	Run: switchAccount,
}

// Predefined configurations
var configs = map[string]Config{
	"naxodev": {
		SSHKey: "naxo-github",
		Name:   "Nacho Vazquez",
		Email:  "nacho@naxo.dev",
	},
	"aster": {
		SSHKey: "github_aster_secure_rsa",
		Name:   "Nacho Vazquez",
		Email:  "nacho@astercare.com",
	},
}

type Config struct {
	SSHKey string
	Name   string
	Email  string
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func switchAccount(cmd *cobra.Command, args []string) {
	var account string
	if len(args) > 0 {
		account = args[0]
	} else {
		account = selectAccount()
	}

	config, exists := configs[account]
	if !exists {
		fmt.Printf("Unknown account: %s\n", account)
		os.Exit(1)
	}

	updateSSHConfig(config.SSHKey)
	updateGitConfig(config.Name, config.Email)

	fmt.Printf("Switched to GitHub account: %s\n", account)
}

func selectAccount() string {
	fmt.Println("Select an account:")
	var options []string
	for account := range configs {
		options = append(options, account)
	}

	for i, opt := range options {
		fmt.Printf("%d. %s\n", i+1, opt)
	}

	var selection int
	for {
		fmt.Print("Enter your choice (1-" + fmt.Sprint(len(options)) + "): ")
		_, err := fmt.Scanf("%d", &selection)
		if err == nil && selection > 0 && selection <= len(options) {
			break
		}
		fmt.Println("Invalid selection. Please try again.")
	}

	return options[selection-1]
}

func updateSSHConfig(sshKey string) {
	sshConfigPath := os.Getenv("HOME") + "/.ssh/config" // Update this path

	input, err := os.ReadFile(sshConfigPath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(input), "\n")
	var newLines []string
	inGithubBlock := false
	githubBlockFound := false

	for _, line := range lines {
		if strings.HasPrefix(line, "Host github.com") {
			inGithubBlock = true
			githubBlockFound = true
			newLines = append(newLines, line)
		} else if inGithubBlock && strings.HasPrefix(line, "Host ") {
			inGithubBlock = false
			newLines = append(newLines, line)
		} else if inGithubBlock && strings.HasPrefix(line, "  IdentityFile ") {
			newLines = append(newLines, fmt.Sprintf("  IdentityFile ~/.ssh/%s", sshKey))
		} else {
			newLines = append(newLines, line)
		}
	}

	if !githubBlockFound {
		newLines = append(newLines, "\nHost github.com")
		newLines = append(newLines, "  AddKeysToAgent yes")
		newLines = append(newLines, "  UseKeychain yes")
		newLines = append(newLines, fmt.Sprintf("  IdentityFile ~/.ssh/%s", sshKey))
	}

	output := strings.Join(newLines, "\n")
	err = os.WriteFile(sshConfigPath, []byte(output), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}
}

func updateGitConfig(name, email string) {
	gitConfigs := map[string]string{
		"user.name":  name,
		"user.email": email,
	}

	for key, value := range gitConfigs {
		err := exec.Command("git", "config", "--global", key, value).Run()
		if err != nil {
			fmt.Printf("Error updating Git config %s: %v\n", key, err)
			os.Exit(1)
		}
	}
}
