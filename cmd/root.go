package cmd

import (
	"context"

	ufcli "github.com/urfave/cli/v3"

	"github.com/mholtzscher/today/cmd/add"
	"github.com/mholtzscher/today/cmd/show"
	"github.com/mholtzscher/today/internal/cli"
)

//nolint:gochecknoglobals // Required for release-please versioning
var Version = "0.1.2" // x-release-please-version

func Run(ctx context.Context, args []string) error {
	app := &ufcli.Command{
		Name:    "today",
		Usage:   "Track daily wins and accomplishments",
		Version: Version,
		Flags: []ufcli.Flag{
			&ufcli.StringFlag{
				Name:    cli.FlagDB,
				Usage:   "Database path",
				Value:   cli.DefaultDBPath(),
				Sources: ufcli.EnvVars("TODAY_DB"),
			},
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
			add.NewCommand(),
			show.NewCommand(),
		},
	}

	return app.Run(ctx, args)
}
