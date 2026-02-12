package main

import (
	"fmt"

	"to/pkg/config"

	"github.com/spf13/cobra"
)

var unregCmd = &cobra.Command{
	Use:           "unreg <alias>",
	Short:         "Unregister an alias",
	Long:          "Remove an alias from the database.",
	Args:          cobra.ExactArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runUnreg,
}

func init() {
	rootCmd.AddCommand(unregCmd)
}

func runUnreg(cmd *cobra.Command, args []string) error {
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
