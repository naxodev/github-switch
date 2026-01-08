// Package config handles loading and saving account configurations.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

type Account struct {
	SSHKey string `yaml:"ssh_key"`
	Name   string `yaml:"name"`
	Email  string `yaml:"email"`
}

type Config struct {
	Accounts map[string]Account `yaml:"accounts"`
}

var configPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		configPath = ".github-switch.yaml"
		return
	}
	configPath = filepath.Join(home, ".github-switch.yaml")
}

func GetConfigPath() string {
	return configPath
}

func Load() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Accounts: make(map[string]Account)}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Accounts == nil {
		cfg.Accounts = make(map[string]Account)
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (c *Config) GetAccount(name string) (Account, bool) {
	acc, ok := c.Accounts[name]
	return acc, ok
}

func (c *Config) AddAccount(name string, acc Account) {
	c.Accounts[name] = acc
}

func (c *Config) RemoveAccount(name string) bool {
	if _, ok := c.Accounts[name]; !ok {
		return false
	}
	delete(c.Accounts, name)
	return true
}

func (c *Config) ListAccounts() []string {
	names := make([]string, 0, len(c.Accounts))
	for name := range c.Accounts {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
