package main

import (
	"io"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// detectTerminal is injected in tests so TTY behavior can be verified
// deterministically without depending on the test runner environment.
var detectTerminal = term.IsTerminal

// fdWriter captures writers that expose an OS file descriptor.
// *os.File satisfies this interface, which lets us ask whether the
// destination stream is an interactive terminal.
type fdWriter interface {
	Fd() uintptr
}

// cliFormatter centralizes all command-layer styling decisions. This keeps
// style calls consistent and lets us cleanly no-op styles for non-TTY output.
type cliFormatter struct {
	enabled bool

	helpLabel *color.Color
	helpText  *color.Color
	example   *color.Color

	listLabel *color.Color
	listAlias *color.Color
	listPath  *color.Color

	success *color.Color
	warn    *color.Color
	err     *color.Color
}

// newCLIFormatterForWriter enables ANSI colors only when the target writer
// appears to be a terminal.
func newCLIFormatterForWriter(w io.Writer) *cliFormatter {
	return newCLIFormatter(isTTYWriter(w))
}

// newCLIFormatter builds the shared style palette for a specific enablement
// mode so callers can choose deterministic behavior in tests.
func newCLIFormatter(enabled bool) *cliFormatter {
	// Build each style through a helper so we can force-enable ANSI output for
	// TTY mode without relying on global color package state.
	makeStyle := func(attrs ...color.Attribute) *color.Color {
		st := color.New(attrs...)
		if enabled {
			st.EnableColor()
		} else {
			st.DisableColor()
		}
		return st
	}

	return &cliFormatter{
		enabled: enabled,

		helpLabel: makeStyle(color.FgCyan, color.Bold),
		helpText:  makeStyle(color.FgWhite),
		example:   makeStyle(color.FgWhite),

		listLabel: makeStyle(color.FgCyan, color.Bold),
		listAlias: makeStyle(color.FgYellow, color.Bold),
		listPath:  makeStyle(color.FgGreen),

		success: makeStyle(color.FgGreen),
		warn:    makeStyle(color.FgYellow),
		err:     makeStyle(color.FgRed, color.Bold),
	}
}

// isTTYWriter reports whether the destination is an interactive terminal.
// Non-file writers (buffers, pipes wrapped without Fd) are treated as non-TTY.
func isTTYWriter(w io.Writer) bool {
	fd, ok := w.(fdWriter)
	if !ok {
		return false
	}
	return detectTerminal(int(fd.Fd()))
}

// style applies a color style only when enabled; otherwise it returns plain
// text so redirected output and snapshots stay stable.
func (f *cliFormatter) style(c *color.Color, text string) string {
	if !f.enabled {
		return text
	}
	return c.Sprint(text)
}

func (f *cliFormatter) HelpLabel(text string) string {
	return f.style(f.helpLabel, text)
}

func (f *cliFormatter) HelpText(text string) string {
	return f.style(f.helpText, text)
}

func (f *cliFormatter) Example(text string) string {
	return f.style(f.example, text)
}

func (f *cliFormatter) ListLabel(text string) string {
	return f.style(f.listLabel, text)
}

func (f *cliFormatter) ListAlias(text string) string {
	return f.style(f.listAlias, text)
}

func (f *cliFormatter) ListPath(text string) string {
	return f.style(f.listPath, text)
}

func (f *cliFormatter) Success(text string) string {
	return f.style(f.success, text)
}

func (f *cliFormatter) Warning(text string) string {
	return f.style(f.warn, text)
}

func (f *cliFormatter) Error(text string) string {
	return f.style(f.err, text)
}

// formatterForOutput resolves styles using the command output stream.
func formatterForOutput(cmd *cobra.Command) *cliFormatter {
	return newCLIFormatterForWriter(cmd.OutOrStdout())
}

// formatterForError resolves styles using the command error stream so warning
// and error coloring follows the actual destination.
func formatterForError(cmd *cobra.Command) *cliFormatter {
	return newCLIFormatterForWriter(cmd.ErrOrStderr())
}
