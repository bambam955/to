package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// flagReg and flagUnreg track which operation mode was requested via
// --reg/-r or --unreg/-u flags. At most one should be true; when neither
// is set the command defaults to navigation mode.
var (
	flagReg   bool
	flagUnreg bool
)

// rootCmd is the top-level cobra command. It uses flags (not subcommands)
// so that the shell wrapper can expose a natural syntax:
//
//	to <alias>              — navigate
//	to --reg <alias> <dir>  — register
//	to --unreg <alias>      — unregister
var rootCmd = &cobra.Command{
	Use:           "to [alias]",
	Short:         "A modern directory navigation tool",
	Long:          "A modern directory navigation tool with JSON database support.",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          run,
}

func init() {
	rootCmd.Flags().BoolVarP(&flagReg, "reg", "r", false, "Register a new alias: to --reg <alias> <directory>")
	rootCmd.Flags().BoolVarP(&flagUnreg, "unreg", "u", false, "Unregister an alias: to --unreg <alias>")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// run dispatches to the appropriate handler based on which flag is set.
// With no flags it defaults to navigation (alias lookup + cd protocol).
func run(cmd *cobra.Command, args []string) error {
	switch {
	case flagReg:
		return runReg(cmd, args)
	case flagUnreg:
		return runUnreg(cmd, args)
	default:
		return runNavigate(cmd, args)
	}
}

// Execute runs the root command. On error it prints a user-facing message
// to stderr and exits with code 1.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
