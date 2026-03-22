package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"to/pkg/config"
)

func runClean(cmd *cobra.Command, args []string) error {
	formatter := formatterForOutput(cmd)

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
		fmt.Fprintln(cmd.OutOrStdout(), formatter.ListLabel("no invalid aliases found"))
		return nil
	}

	if err := db.Save(dbPath); err != nil {
		return err
	}

	for _, a := range removed {
		fmt.Fprintln(
			cmd.OutOrStdout(),
			formatter.Warning(fmt.Sprintf("removed '%s' -> %s", a.Name, a.Directory)),
		)
	}

	summary := fmt.Sprintf("cleaned %d alias(es), %d remaining", len(removed), len(db.ListAliases()))
	fmt.Fprintln(cmd.OutOrStdout(), formatter.Success(summary))
	return nil
}
