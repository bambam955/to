package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"to/pkg/config"
)

func runExp(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: to --exp <alias>")
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

	entry, err := db.GetAlias(name)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), entry.Directory)
	return nil
}
