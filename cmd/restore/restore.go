package restore

import (
	"context"
	"errors"
	"fmt"

	ufcli "github.com/urfave/cli/v3"

	"github.com/mholtzscher/today/internal/cli"
	"github.com/mholtzscher/today/internal/db"
	"github.com/mholtzscher/today/internal/output"
)

func NewCommand() *ufcli.Command {
	return &ufcli.Command{
		Name:  "restore",
		Usage: "Restore an archived entry by id",
		Arguments: []ufcli.Argument{
			&ufcli.IntArg{Name: "id"},
		},
		Action: func(ctx context.Context, cmd *ufcli.Command) error {
			id := cmd.IntArg("id")
			if id <= 0 {
				return errors.New("id argument required")
			}

			database, err := db.Open(cli.GetDBPath(cmd))
			if err != nil {
				return err
			}
			defer database.Close()

			store := db.NewStore(database)

			entry, err := store.GetEntry(ctx, int64(id))
			if err != nil {
				if errors.Is(err, db.ErrEntryNotFound) {
					output.Stdoutln("No entry restored")
					return nil
				}
				return err
			}

			if entry.ArchivedAt == nil {
				output.Stdoutln("No entry restored")
				return nil
			}

			restored, err := store.RestoreEntry(ctx, int64(id))
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
