# github-switch

A CLI tool to quickly switch between multiple GitHub accounts by updating SSH config and Git configuration.

## Installation

```bash
go install github.com/naxodev/github-switch@latest
```

Or build from source:

```bash
git clone https://github.com/naxodev/github-switch.git
cd github-switch
make install
```

## Quick Start

1. Initialize the config file:

```bash
github-switch init
```

2. Add your GitHub accounts:

```bash
github-switch add personal
# Follow the prompts to enter your name, email, and SSH key

github-switch add work --name "Your Name" --email "you@company.com" --ssh-key "id_work_rsa"
```

3. Switch between accounts:

```bash
github-switch switch personal
# or use the shorthand
github-switch sw work
```

## Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `switch <account>` | `sw` | Switch to a GitHub account |
| `list` | `ls` | List all configured accounts |
| `current` | | Show current Git/SSH configuration |
| `add <name>` | | Add a new account |
| `remove <name>` | `rm` | Remove an account |
| `init` | | Initialize config file |

## Configuration

Accounts are stored in `~/.github-switch.yaml`:

```yaml
accounts:
  personal:
    ssh_key: id_personal_rsa
    name: Your Name
    email: you@example.com
  work:
    ssh_key: id_work_rsa
    name: Your Name
    email: you@company.com
```

## What It Does

When you switch accounts, `github-switch`:

1. Updates `~/.ssh/config` to use the correct SSH key for `github.com`
2. Sets global Git `user.name` and `user.email`
3. Adds the SSH key to your ssh-agent

## Prerequisites

- Go 1.21+ (for installation)
- SSH keys configured for each GitHub account
- Git installed

## License

MIT License - see [LICENSE](LICENSE) for details.
