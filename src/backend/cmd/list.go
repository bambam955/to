package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"to/pkg/config"
)

func runList(cmd *cobra.Command, args []string) error {
	formatter := formatterForOutput(cmd)

	dbPath, err := config.GetDatabasePath()
	if err != nil {
		return err
	}

	db, err := loadOrInitDB(dbPath)
	if err != nil {
		return err
	}

	aliases := db.ListAliases()
	if len(aliases) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), formatter.ListLabel("no aliases registered"))
		return nil
	}

	// Sort alphabetically by name.
	sort.Slice(aliases, func(i, j int) bool {
		return aliases[i].Name < aliases[j].Name
	})

	// Find the longest alias name for column alignment.
	maxLen := 0
	for _, a := range aliases {
		if len(a.Name) > maxLen {
			maxLen = len(a.Name)
		}
	}

	for _, a := range aliases {
		// Pad before styling so visual alignment remains stable in terminals.
		aliasColumn := fmt.Sprintf("%-*s", maxLen, a.Name)
		fmt.Fprintf(
			cmd.OutOrStdout(),
			"%s  %s\n",
			formatter.ListAlias(aliasColumn),
			formatter.ListPath(a.Directory),
		)
	}

	return nil
}
