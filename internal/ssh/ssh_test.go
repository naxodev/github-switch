package ssh

import (
	"testing"
)

func TestUpdateGitHubBlock(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		sshKey   string
		expected string
	}{
		{
			name:   "update existing github block",
			sshKey: "new-key",
			input: `Host github.com
  AddKeysToAgent yes
  IdentityFile ~/.ssh/old-key

Host gitlab.com
  IdentityFile ~/.ssh/gitlab-key`,
			expected: `Host github.com
  AddKeysToAgent yes
  IdentityFile ~/.ssh/new-key

Host gitlab.com
  IdentityFile ~/.ssh/gitlab-key`,
		},
		{
			name:   "add github block when missing",
			sshKey: "my-key",
			input: `Host gitlab.com
  IdentityFile ~/.ssh/gitlab-key`,
			expected: `Host gitlab.com
  IdentityFile ~/.ssh/gitlab-key

Host github.com
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/my-key`,
		},
		{
			name:     "create from empty",
			sshKey:   "test-key",
			input:    "",
			expected: "\nHost github.com\n  AddKeysToAgent yes\n  UseKeychain yes\n  IdentityFile ~/.ssh/test-key",
		},
		{
			name:   "preserve indentation with tabs",
			sshKey: "new-key",
			input: `Host github.com
	AddKeysToAgent yes
	IdentityFile ~/.ssh/old-key`,
			expected: `Host github.com
	AddKeysToAgent yes
	IdentityFile ~/.ssh/new-key`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := updateGitHubBlock(tt.input, tt.sshKey)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("mismatch:\nexpected:\n%s\n\ngot:\n%s", tt.expected, result)
			}
		})
	}
}
