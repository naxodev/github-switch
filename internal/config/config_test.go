package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test-config.yaml")
	configPath = testPath

	cfg := &Config{
		Accounts: map[string]Account{
			"test": {
				SSHKey: "test-key",
				Name:   "Test User",
				Email:  "test@example.com",
			},
		},
	}

	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	acc, ok := loaded.GetAccount("test")
	if !ok {
		t.Fatal("expected account 'test' to exist")
	}

	if acc.SSHKey != "test-key" {
		t.Errorf("expected ssh_key 'test-key', got '%s'", acc.SSHKey)
	}
	if acc.Name != "Test User" {
		t.Errorf("expected name 'Test User', got '%s'", acc.Name)
	}
	if acc.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", acc.Email)
	}
}

func TestListAccounts(t *testing.T) {
	cfg := &Config{
		Accounts: map[string]Account{
			"charlie": {},
			"alpha":   {},
			"bravo":   {},
		},
	}

	accounts := cfg.ListAccounts()
	expected := []string{"alpha", "bravo", "charlie"}

	if len(accounts) != len(expected) {
		t.Fatalf("expected %d accounts, got %d", len(expected), len(accounts))
	}

	for i, name := range accounts {
		if name != expected[i] {
			t.Errorf("expected account[%d] = '%s', got '%s'", i, expected[i], name)
		}
	}
}

func TestRemoveAccount(t *testing.T) {
	cfg := &Config{
		Accounts: map[string]Account{
			"test": {},
		},
	}

	if !cfg.RemoveAccount("test") {
		t.Error("expected RemoveAccount to return true")
	}

	if cfg.RemoveAccount("nonexistent") {
		t.Error("expected RemoveAccount to return false for nonexistent account")
	}

	if _, ok := cfg.GetAccount("test"); ok {
		t.Error("expected account 'test' to be removed")
	}
}

func TestLoadNonexistentConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath = filepath.Join(tmpDir, "nonexistent.yaml")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Accounts == nil {
		t.Error("expected Accounts map to be initialized")
	}

	if len(cfg.Accounts) != 0 {
		t.Error("expected empty Accounts map")
	}
}

func TestConfigPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test-config.yaml")
	configPath = testPath

	cfg := &Config{Accounts: map[string]Account{}}
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	info, err := os.Stat(testPath)
	if err != nil {
		t.Fatalf("failed to stat config file: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0o600 {
		t.Errorf("expected permissions 0600, got %o", perm)
	}
}
