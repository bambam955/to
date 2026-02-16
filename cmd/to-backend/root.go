package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Operation-mode flags select which action to perform. At most one should
// be true; when none is set the command defaults to navigation mode.
var (
	flagReg   bool
	flagUnreg bool
	flagList  bool
	flagClean bool
	flagExp   bool
)

// rootCmd is the top-level cobra command. It uses flags (not subcommands)
// so that the shell wrapper can expose a natural syntax:
//
//	to <alias>              — navigate
//	to --reg <alias> <dir>  — register
//	to --unreg <alias>      — unregister
//	to --list               — list all aliases
//	to --clean              — remove invalid aliases
//	to --exp <alias>        — expand alias to path
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
	rootCmd.Flags().BoolVarP(&flagList, "list", "l", false, "List all registered aliases")
	rootCmd.Flags().BoolVarP(&flagClean, "clean", "c", false, "Remove aliases pointing to directories that no longer exist")
	rootCmd.Flags().BoolVarP(&flagExp, "exp", "e", false, "Show the full directory path for an alias: to --exp <alias>")
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
	case flagList:
		return runList(cmd, args)
	case flagClean:
		return runClean(cmd, args)
	case flagExp:
		return runExp(cmd, args)
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
