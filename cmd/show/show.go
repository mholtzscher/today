package show

import (
	"context"
	"fmt"

	ufcli "github.com/urfave/cli/v3"

	"github.com/mholtzscher/today/internal/cli"
	"github.com/mholtzscher/today/internal/db"
	"github.com/mholtzscher/today/internal/output"
)

func NewCommand() *ufcli.Command {
	return &ufcli.Command{
		Name:    "show",
		Usage:   "Show entries",
		Aliases: []string{"s"},
		Flags: []ufcli.Flag{
			&ufcli.IntFlag{
				Name:  "days",
				Value: 0,
				Usage: "Number of days to show (0 = today only)",
			},
			&ufcli.BoolFlag{
				Name:  "all",
				Usage: "Include archived entries",
			},
		},
		Arguments: []ufcli.Argument{
			&ufcli.IntArg{
				Name: "days-arg",
			},
		},
		Action: func(ctx context.Context, cmd *ufcli.Command) error {
			days := cmd.Int("days")
			if days == 0 {
				daysArg := cmd.IntArg("days-arg")
				if daysArg > 0 {
					days = daysArg
				}
			}
			if days == 0 {
				days = 1
			}

			database, err := db.Open(cli.GetDBPath(cmd))
			if err != nil {
				return err
			}
			defer database.Close()

			store := db.NewStore(database)
			entries, err := store.ListEntries(ctx, days, cmd.Bool("all"))
			if err != nil {
				return err
			}

			printEntries(entries)
			return nil
		},
	}
}

func printEntries(entries []db.Entry) {
	if len(entries) == 0 {
		output.Stdoutln("No entries found")
		return
	}

	var currentDate string
	for _, e := range entries {
		date := e.CreatedAt.Format("2006-01-02")
		if date != currentDate {
			if currentDate != "" {
				output.Stdoutln("")
			}
			currentDate = date
			output.Stdoutln(fmt.Sprintf("=== %s ===", date))
		}

		line := fmt.Sprintf("• #%d %s", e.ID, e.Text)
		if e.ArchivedAt != nil {
			line = fmt.Sprintf("• #%d [archived] %s", e.ID, e.Text)
		}
		output.Stdoutln(line)
	}
}
