package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"to/pkg/config"
	"to/pkg/database"
	"to/pkg/protocol"
)

const controlChannelFD = 3

// writeNavigationControl emits navigation control frames to the wrapper-owned
// control channel. Tests replace this to assert navigation behavior without
// needing real file-descriptor plumbing.
var writeNavigationControl = func(path string) error {
	control := os.NewFile(uintptr(controlChannelFD), "to-control-channel")
	if control == nil {
		return fmt.Errorf("control channel unavailable")
	}
	defer control.Close()

	if err := protocol.WriteNavigationControlFrame(control, path); err != nil {
		return fmt.Errorf("control channel unavailable: %w", err)
	}
	return nil
}

// runNavigate resolves an alias to its directory, verifies the directory
// still exists, updates the last-visited timestamp, and emits a protocol
// control response on fd 3 that the shell wrapper parses to perform cd.
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

	if err := writeNavigationControl(entry.Directory); err != nil {
		return fmt.Errorf("navigation protocol channel unavailable")
	}

	return nil
}
