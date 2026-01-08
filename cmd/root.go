// Package cmd implements the CLI commands for github-switch.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "github-switch",
	Short: "Switch between GitHub accounts",
	Long: `A CLI tool to switch between different GitHub accounts
by modifying SSH config and Git configuration.

Use 'github-switch switch <account>' to switch accounts,
or 'github-switch list' to see available accounts.`,
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("github-switch version {{.Version}}\n")
}
