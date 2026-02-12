package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"to/pkg/config"
	"to/pkg/database"
	"to/pkg/errors"
)

var regCmd = &cobra.Command{
	Use:           "reg <alias> <directory>",
	Short:         "Register a new alias",
	Long:          "Register a new alias pointing to a directory.",
	Args:          cobra.ExactArgs(2),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runReg,
}

func init() {
	rootCmd.AddCommand(regCmd)
}

// loadOrInitDB loads the database from the given path, or initializes a new
// one if the file does not exist.
func loadOrInitDB(dbPath string) (*database.Database, error) {
	db, err := database.Load(dbPath)
	if err != nil {
		// If the database file doesn't exist, initialize a new one.
		if e, ok := err.(*errors.Error); ok && e.Type == errors.ErrorTypeNotFound {
			return database.InitDatabase(dbPath)
		}
		return nil, err
	}
	return db, nil
}

// reservedNames returns the set of subcommand names that cannot be used as aliases.
func reservedNames() map[string]bool {
	reserved := make(map[string]bool)
	for _, cmd := range rootCmd.Commands() {
		reserved[cmd.Name()] = true
	}
	return reserved
}

func runReg(cmd *cobra.Command, args []string) error {
	name := args[0]
	dir := args[1]

	// Reject alias names that conflict with subcommand names.
	if reservedNames()[name] {
		return fmt.Errorf("'%s' is a reserved command name", name)
	}

	// Resolve to absolute path.
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	// Get database path.
	dbPath, err := config.GetDatabasePath()
	if err != nil {
		return err
	}

	// Load or initialize the database.
	db, err := loadOrInitDB(dbPath)
	if err != nil {
		return err
	}

	// Check for duplicate directories and warn.
	existing := db.FindByDirectory(absDir)
	if len(existing) > 0 {
		names := make([]string, len(existing))
		for i, a := range existing {
			names[i] = a.Name
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "warning: directory already registered as: %s\n", strings.Join(names, ", "))
	}

	// Add the alias (validates name, directory, and checks duplicate names).
	if err := db.AddAlias(name, absDir); err != nil {
		return err
	}

	// Save the database.
	if err := db.Save(dbPath); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "registered '%s' -> %s\n", name, absDir)
	return nil
}
