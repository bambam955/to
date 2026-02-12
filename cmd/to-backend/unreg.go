package main

import (
	"fmt"

	"to/pkg/config"

	"github.com/spf13/cobra"
)

// runUnreg removes an alias from the database. Returns an error if the
// alias does not exist.
func runUnreg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: to --unreg <alias>")
	}

	name := args[0]

	dbPath, err := config.GetDatabasePath()
	if err != nil {
		return err
	}

	db, err := loadOrInitDB(dbPath)
	if err != nil {
		return err
	}

	if err := db.RemoveAlias(name); err != nil {
		return err
	}

	if err := db.Save(dbPath); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "unregistered '%s'\n", name)
	return nil
}
