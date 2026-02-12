package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"to/pkg/config"
	"to/pkg/database"
	"to/pkg/protocol"
)

func init() {
	rootCmd.Use = "to [alias]"
	rootCmd.Args = cobra.ExactArgs(1)
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.RunE = runNavigate
}

func runNavigate(cmd *cobra.Command, args []string) error {
	alias := args[0]

	dbPath, err := config.GetDatabasePath()
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "error: %s\n", err)
		return err
	}

	db, err := loadOrInitDB(dbPath)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "error: %s\n", err)
		return err
	}

	entry, err := db.GetAlias(alias)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "error: %s\n", err)
		return err
	}

	if err := database.ValidateDirectory(entry.Directory); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "error: directory no longer exists: %s\n", entry.Directory)
		return err
	}

	if err := db.UpdateLastVisited(alias); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "error: %s\n", err)
		return err
	}

	if err := db.Save(dbPath); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "error: %s\n", err)
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), protocol.NavigationResponse(entry.Directory))
	return nil
}
