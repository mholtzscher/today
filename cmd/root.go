// Package cmd implements the CLI commands for today.
package cmd

import (
	"context"

	ufcli "github.com/urfave/cli/v3"
	"github.com/mholtzscher/today/cmd/example"
	"github.com/mholtzscher/today/internal/cli"
)

// Version is set at build time.
var Version = "0.1.0" // x-release-please-version

// Run is the entry point for the CLI.
func Run(ctx context.Context, args []string) error {
	app := &ufcli.Command{
		Name:    "today",
		Usage:   "A Go CLI tool built with Nix",
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
			example.NewCommand(),
		},
	}

	return app.Run(ctx, args)
}
