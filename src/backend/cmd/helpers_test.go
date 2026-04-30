package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

// resetFlags resets all root command flags to their default values.
// Must be called before each test that uses rootCmd.Execute().
func resetFlags(t *testing.T) {
	t.Helper()

	previousBuildVersion := buildVersion

	resetCobraFlags := func() {
		// Cobra tracks flag values in the FlagSet in addition to the package
		// bool variables, so we reset both to avoid state leakage across tests.
		rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
			_ = f.Value.Set(f.DefValue)
			f.Changed = false
		})
	}

	t.Cleanup(func() {
		buildVersion = previousBuildVersion
		configureVersion(rootCmd)
		resetCobraFlags()
		flagReg = false
		flagUnreg = false
		flagList = false
		flagClean = false
		flagExp = false
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(io.Discard)
		rootCmd.SetErr(io.Discard)
		flagUninstall = false
		flagPurge = false
	})
	flagReg = false
	flagUnreg = false
	flagList = false
	flagClean = false
	flagExp = false
	flagUninstall = false
	flagPurge = false
	configureVersion(rootCmd)
	resetCobraFlags()
	rootCmd.SetArgs(nil)
	rootCmd.SetOut(io.Discard)
	rootCmd.SetErr(io.Discard)
}

// setBuildVersion updates the shared build metadata for tests and restores the
// previous value after the test finishes.
func setBuildVersion(t *testing.T, version string) {
	t.Helper()

	previousBuildVersion := buildVersion
	buildVersion = version
	configureVersion(rootCmd)
	t.Cleanup(func() {
		buildVersion = previousBuildVersion
		configureVersion(rootCmd)
	})
}

// testTTYBuffer captures output while exposing a fake file descriptor so TTY
// detection can be controlled in tests.
type testTTYBuffer struct {
	bytes.Buffer
	fd uintptr
}

func (b *testTTYBuffer) Fd() uintptr {
	return b.fd
}

// forceTTYDetection overrides terminal detection for tests that need explicit
// TTY/non-TTY behavior independent from the environment.
func forceTTYDetection(t *testing.T, isTTY bool) {
	t.Helper()

	previous := detectTerminal
	detectTerminal = func(fd int) bool {
		return isTTY
	}
	t.Cleanup(func() {
		detectTerminal = previous
	})
}

func containsANSI(text string) bool {
	return strings.Contains(text, "\x1b[")
}

// stubNavigationControlWriter swaps the fd3 navigation writer for tests so
// command behavior can be verified without touching real process descriptors.
func stubNavigationControlWriter(t *testing.T, fn func(path string) error) {
	t.Helper()

	previous := writeNavigationControl
	writeNavigationControl = fn
	t.Cleanup(func() {
		writeNavigationControl = previous
	})
}
