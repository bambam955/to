package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"to/pkg/protocol"
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

// rootExamples documents the common root-command workflows in a compact
// command-plus-description layout for `to --help`.
const rootExamples = `  $ to work                         Navigate to the "work" alias
  $ to --reg work ~/projects/work    Register "work" for a directory
  $ to --list                        List registered aliases
  $ to --exp work                    Print the path for "work"
  $ to --unreg work                  Remove the "work" alias
  $ to --clean                       Remove aliases for missing directories`

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
	Example:       rootExamples,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          run,
}

// styledHelpTemplate keeps Cobra help output readable in interactive terminals
// while preserving plain output for non-TTY destinations.
const styledHelpTemplate = `{{with (or .Long .Short)}}{{toHelpText $ .}}{{"\n\n"}}{{end}}` +
	`{{if or .Runnable .HasSubCommands}}{{toHelpLabel . "Usage:"}}{{"\n  "}}{{toHelpText . .UseLine}}{{"\n"}}{{end}}` +
	`{{if .HasAvailableSubCommands}}{{"\n"}}{{toHelpLabel . "Available Commands:"}}{{"\n"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}{{"  "}}{{toHelpText $ (rpad .Name .NamePadding)}} {{toHelpText $ .Short}}{{"\n"}}{{end}}{{end}}{{end}}` +
	`{{if .HasAvailableLocalFlags}}{{"\n"}}{{toHelpLabel . "Flags:"}}{{"\n"}}{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}` +
	`{{if .HasAvailableInheritedFlags}}{{"\n"}}{{toHelpLabel . "Global Flags:"}}{{"\n"}}{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}` +
	`{{if .HasExample}}{{"\n\n"}}{{toHelpLabel . "Examples:"}}{{"\n"}}{{toHelpExample . .Example}}{{"\n"}}{{end}}` +
	`{{if .HasHelpSubCommands}}{{"\n"}}{{toHelpLabel . "Additional help topics:"}}{{"\n"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}{{"  "}}{{toHelpText $ (rpad .CommandPath .CommandPathPadding)}} {{toHelpText $ .Short}}{{"\n"}}{{end}}{{end}}{{end}}` +
	`{{if .HasAvailableSubCommands}}{{"\n"}}{{toHelpText . (printf "Use \"%s [command] --help\" for more information about a command." .CommandPath)}}{{end}}{{"\n"}}`

func init() {
	configureVersion(rootCmd)

	// Register template functions once; each call resolves styles from the
	// command's active output stream so TTY/non-TTY behavior stays correct.
	cobra.AddTemplateFunc("toHelpLabel", func(cmd *cobra.Command, text string) string {
		return formatterForOutput(cmd).HelpLabel(text)
	})
	cobra.AddTemplateFunc("toHelpText", func(cmd *cobra.Command, text string) string {
		return formatterForOutput(cmd).HelpText(text)
	})
	cobra.AddTemplateFunc("toHelpExample", func(cmd *cobra.Command, text string) string {
		lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
		formatter := formatterForOutput(cmd)
		for i := range lines {
			lines[i] = formatter.Example(lines[i])
		}
		return strings.Join(lines, "\n")
	})

	rootCmd.SetHelpTemplate(styledHelpTemplate)

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
		formatter := newCLIFormatterForWriter(os.Stderr)
		fmt.Fprintln(os.Stderr, formatCommandError(err, formatter))
		os.Exit(1)
	}
}

// formatCommandError ensures all command-level errors keep the same
// "error: ..." shape while applying terminal styling when enabled.
func formatCommandError(err error, formatter *cliFormatter) string {
	return formatter.Error(protocol.ErrorResponse(err.Error()))
}
