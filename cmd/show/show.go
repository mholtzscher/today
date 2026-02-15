package show

import (
	"context"
	"fmt"
	"io"
	"os"

	ufcli "github.com/urfave/cli/v3"

	"github.com/mholtzscher/today/internal/cli"
	"github.com/mholtzscher/today/internal/db"
	"github.com/mholtzscher/today/internal/entry"
)

func NewCommand() *ufcli.Command {
	return &ufcli.Command{
		Name:  "show",
		Usage: "Show entries",
		Flags: []ufcli.Flag{
			&ufcli.IntFlag{
				Name:  "days",
				Value: 0,
				Usage: "Number of days to show (0 = today only)",
			},
		},
		Arguments: []ufcli.Argument{
			&ufcli.IntArg{
				Name: "days-arg",
			},
		},
		Action: func(_ context.Context, cmd *ufcli.Command) error {
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

			store := entry.NewStore(database)
			entries, err := store.GetByDays(days)
			if err != nil {
				return err
			}

			printEntries(os.Stdout, entries)
			return nil
		},
	}
}

func printEntries(w io.Writer, entries []entry.Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No entries found")
		return
	}

	var currentDate string
	for _, e := range entries {
		date := e.CreatedAt.Format("2006-01-02")
		if date != currentDate {
			currentDate = date
			fmt.Fprintf(w, "\n=== %s ===\n", date)
		}
		fmt.Fprintf(w, "• %s\n", e.Text)
	}
}
