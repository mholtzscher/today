// Package cmd implements the CLI commands for today.
package cmd

import (
	"context"
	"os"
	"path/filepath"

	ufcli "github.com/urfave/cli/v3"

	"github.com/mholtzscher/today/cmd/add"
	"github.com/mholtzscher/today/cmd/show"
	"github.com/mholtzscher/today/internal/cli"
)

// Version is set at build time.
//
//nolint:gochecknoglobals // Required for release-please versioning
var Version = "0.1.0" // x-release-please-version

func defaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "today.db"
	}
	return filepath.Join(home, "today.db")
}

// Run is the entry point for the CLI.
func Run(ctx context.Context, args []string) error {
	dbPath := defaultDBPath()

	app := &ufcli.Command{
		Name:    "today",
		Usage:   "Track daily wins and accomplishments",
		Version: Version,
		Flags: []ufcli.Flag{
			&ufcli.BoolFlag{
				Name:  cli.FlagVerbose,
				Usage: "Print verbose output",
			},
			&ufcli.BoolFlag{
				Name:  cli.FlagNoColor,
				Usage: "Disable colored output",
			},
		},
		Commands: []*ufcli.Command{
			add.NewCommand(dbPath),
			show.NewCommand(dbPath),
		},
	}

	return app.Run(ctx, args)
}
