// Package ssh handles SSH config file operations.
package ssh

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".ssh", "config"), nil
}

func UpdateConfig(sshKey string) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	input, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return createNewConfig(configPath, sshKey)
		}
		return fmt.Errorf("failed to read SSH config: %w", err)
	}

	output, err := updateGitHubBlock(string(input), sshKey)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, []byte(output), 0o600); err != nil {
		return fmt.Errorf("failed to write SSH config: %w", err)
	}

	return nil
}

func createNewConfig(path, sshKey string) error {
	content := fmt.Sprintf(`Host github.com
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/%s
`, sshKey)

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("failed to write SSH config: %w", err)
	}

	return nil
}

func updateGitHubBlock(input, sshKey string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var lines []string
	inGithubBlock := false
	githubBlockFound := false
	identityUpdated := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.EqualFold(trimmed, "Host github.com") {
			inGithubBlock = true
			githubBlockFound = true
			lines = append(lines, line)
			continue
		}

		if inGithubBlock && strings.HasPrefix(trimmed, "Host ") {
			inGithubBlock = false
		}

		if inGithubBlock && strings.HasPrefix(trimmed, "IdentityFile") {
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			lines = append(lines, fmt.Sprintf("%sIdentityFile ~/.ssh/%s", indent, sshKey))
			identityUpdated = true
			continue
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to parse SSH config: %w", err)
	}

	if !githubBlockFound {
		lines = append(lines, "", "Host github.com")
		lines = append(lines, "  AddKeysToAgent yes")
		lines = append(lines, "  UseKeychain yes")
		lines = append(lines, fmt.Sprintf("  IdentityFile ~/.ssh/%s", sshKey))
	} else if !identityUpdated {
		for i, line := range lines {
			if strings.EqualFold(strings.TrimSpace(line), "Host github.com") {
				insertLines := []string{fmt.Sprintf("  IdentityFile ~/.ssh/%s", sshKey)}
				lines = append(lines[:i+1], append(insertLines, lines[i+1:]...)...)
				break
			}
		}
	}

	return strings.Join(lines, "\n"), nil
}

func GetCurrentKey() (string, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return "", err
	}

	input, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read SSH config: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(input)))
	inGithubBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.EqualFold(trimmed, "Host github.com") {
			inGithubBlock = true
			continue
		}

		if inGithubBlock && strings.HasPrefix(trimmed, "Host ") {
			break
		}

		if inGithubBlock && strings.HasPrefix(trimmed, "IdentityFile") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				return filepath.Base(parts[1]), nil
			}
		}
	}

	return "", nil
}

func AddKeyToAgent(sshKey string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	keyPath := filepath.Join(home, ".ssh", sshKey)
	cmd := exec.Command("ssh-add", keyPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add key to ssh-agent: %w", err)
	}

	return nil
}
