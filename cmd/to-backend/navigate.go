package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"to/pkg/config"
	"to/pkg/database"
	"to/pkg/protocol"
)

func runNavigate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: to <alias>")
	}

	alias := args[0]

	dbPath, err := config.GetDatabasePath()
	if err != nil {
		return err
	}

	db, err := loadOrInitDB(dbPath)
	if err != nil {
		return err
	}

	entry, err := db.GetAlias(alias)
	if err != nil {
		return err
	}

	if err := database.ValidateDirectory(entry.Directory); err != nil {
		return fmt.Errorf("directory no longer exists: %s", entry.Directory)
	}

	if err := db.UpdateLastVisited(alias); err != nil {
		return err
	}

	if err := db.Save(dbPath); err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), protocol.NavigationResponse(entry.Directory))
	return nil
}
