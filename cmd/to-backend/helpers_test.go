package main

import "testing"

// resetFlags resets all root command flags to their default values.
// Must be called before each test that uses rootCmd.Execute().
func resetFlags(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		flagReg = false
		flagUnreg = false
	})
	flagReg = false
	flagUnreg = false
}
