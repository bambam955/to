package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"to/pkg/config"
)

func runClean(cmd *cobra.Command, args []string) error {
	dbPath, err := config.GetDatabasePath()
	if err != nil {
		return err
	}

	db, err := loadOrInitDB(dbPath)
	if err != nil {
		return err
	}

	removed := db.CleanInvalid()

	if len(removed) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no invalid aliases found")
		return nil
	}

	if err := db.Save(dbPath); err != nil {
		return err
	}

	for _, a := range removed {
		fmt.Fprintf(cmd.OutOrStdout(), "removed '%s' -> %s\n", a.Name, a.Directory)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "cleaned %d alias(es), %d remaining\n", len(removed), len(db.ListAliases()))
	return nil
}
