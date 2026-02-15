package add

import (
	"context"
	"errors"
	"fmt"
	"os"

	ufcli "github.com/urfave/cli/v3"

	"github.com/mholtzscher/today/internal/db"
	"github.com/mholtzscher/today/internal/entry"
)

func NewCommand(dbPath string) *ufcli.Command {
	return &ufcli.Command{
		Name:  "add",
		Usage: "Add a new entry",
		Arguments: []ufcli.Argument{
			&ufcli.StringArg{
				Name: "text",
			},
		},
		Action: func(_ context.Context, cmd *ufcli.Command) error {
			text := cmd.StringArg("text")
			if text == "" {
				return errors.New("text argument required")
			}

			database, err := db.Open(dbPath)
			if err != nil {
				return err
			}
			defer database.Close()

			store := entry.NewStore(database)
			if insertErr := store.Insert(text); insertErr != nil {
				return insertErr
			}

			fmt.Fprintln(os.Stdout, "Added entry")
			return nil
		},
	}
}
