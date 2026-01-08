package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func UpdateGlobalConfig(name, email string) error {
	configs := map[string]string{
		"user.name":  name,
		"user.email": email,
	}

	for key, value := range configs {
		cmd := exec.Command("git", "config", "--global", key, value)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set %s: %w", key, err)
		}
	}

	return nil
}

func GetGlobalConfig(key string) (string, error) {
	cmd := exec.Command("git", "config", "--global", "--get", key)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", fmt.Errorf("failed to get %s: %w", key, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func GetCurrentUser() (name, email string, err error) {
	name, err = GetGlobalConfig("user.name")
	if err != nil {
		return "", "", err
	}

	email, err = GetGlobalConfig("user.email")
	if err != nil {
		return "", "", err
	}

	return name, email, nil
}
