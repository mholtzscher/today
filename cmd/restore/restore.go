package restore

import (
	"context"
	"errors"
	"fmt"

	ufcli "github.com/urfave/cli/v3"

	"github.com/mholtzscher/today/internal/cli"
	"github.com/mholtzscher/today/internal/db"
	"github.com/mholtzscher/today/internal/entry"
	"github.com/mholtzscher/today/internal/output"
)

func NewCommand() *ufcli.Command {
	return &ufcli.Command{
		Name:  "restore",
		Usage: "Restore a soft-deleted entry by id",
		Arguments: []ufcli.Argument{
			&ufcli.IntArg{Name: "id"},
		},
		Action: func(_ context.Context, cmd *ufcli.Command) error {
			id := cmd.IntArg("id")
			if id <= 0 {
				return errors.New("id argument required")
			}

			database, err := db.Open(cli.GetDBPath(cmd))
			if err != nil {
				return err
			}
			defer database.Close()

			store := entry.NewStore(database)
			restored, err := store.RestoreByID(int64(id))
			if err != nil {
				return err
			}

			if !restored {
				output.Stdoutln("No entry restored")
				return nil
			}

			output.Stdoutln(fmt.Sprintf("Restored #%d", id))
			return nil
		},
	}
}
