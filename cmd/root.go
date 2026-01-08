// Package cmd implements the CLI commands for github-switch.
package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "github-switch",
	Short: "Switch between GitHub accounts",
	Long: `A CLI tool to switch between different GitHub accounts
by modifying SSH config and Git configuration.

Use 'github-switch switch <account>' to switch accounts,
or 'github-switch list' to see available accounts.`,
}

func init() {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
			version = info.Main.Version
		}
	}
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("github-switch version {{.Version}}\n")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
