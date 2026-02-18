package add

import (
	"context"
	"errors"

	ufcli "github.com/urfave/cli/v3"

	"github.com/mholtzscher/today/internal/cli"
	"github.com/mholtzscher/today/internal/db"
	"github.com/mholtzscher/today/internal/entry"
	"github.com/mholtzscher/today/internal/output"
)

func NewCommand() *ufcli.Command {
	return &ufcli.Command{
		Name:    "add",
		Usage:   "Add a new entry",
		Aliases: []string{"i", "a"},
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
			svc := entry.NewService(store)
			if insertErr := svc.CreateEntry(ctx, text); insertErr != nil {
				return insertErr
			}

			output.Stdoutln("Added entry")
			return nil
		},
	}
}
