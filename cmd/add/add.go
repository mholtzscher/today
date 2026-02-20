package add

import (
	"context"
	"errors"
	"fmt"
	"time"

	ufcli "github.com/urfave/cli/v3"

	"github.com/mholtzscher/today/internal/cli"
	"github.com/mholtzscher/today/internal/db"
	"github.com/mholtzscher/today/internal/output"
)

func NewCommand() *ufcli.Command {
	return &ufcli.Command{
		Name:    "add",
		Usage:   "Add a new entry",
		Aliases: []string{"i", "a"},
		Flags: []ufcli.Flag{
			&ufcli.StringFlag{
				Name:  "date",
				Usage: "Set entry date (YYYY-MM-DD, past only)",
			},
		},
		Arguments: []ufcli.Argument{
			&ufcli.StringArg{
				Name: "text",
			},
		},
		Action: func(ctx context.Context, cmd *ufcli.Command) error {
			text := cmd.StringArg("text")
			if text == "" {
				return errors.New("text argument required")
			}

			database, err := db.Open(cli.GetDBPath(cmd))
			if err != nil {
				return err
			}
			defer database.Close()

			store := db.NewStore(database)
			dateValue := cmd.String("date")
			if dateValue == "" {
				if insertErr := store.CreateEntry(ctx, text); insertErr != nil {
					return insertErr
				}
			} else {
				createdAt, parseErr := parseEntryDate(dateValue)
				if parseErr != nil {
					return parseErr
				}

				if insertErr := store.CreateEntryAt(ctx, text, createdAt); insertErr != nil {
					return insertErr
				}
			}

			output.Stdoutln("Added entry")
			return nil
		},
	}
}

func parseEntryDate(value string) (time.Time, error) {
	date, err := time.ParseInLocation("2006-01-02", value, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid --date %q, expected YYYY-MM-DD", value)
	}

	now := time.Now().In(time.Local)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if !date.Before(todayStart) {
		return time.Time{}, errors.New("date must be in the past")
	}

	return date, nil
}
