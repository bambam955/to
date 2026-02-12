package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagReg   bool
	flagUnreg bool
)

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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
