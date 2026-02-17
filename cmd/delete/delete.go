package deletecmd

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
		Name:  "delete",
		Usage: "Soft delete an entry by id",
		Flags: []ufcli.Flag{
			&ufcli.BoolFlag{
				Name:  "yes",
				Usage: "Skip confirmation prompt",
			},
		},
		Arguments: []ufcli.Argument{
			&ufcli.IntArg{Name: "id"},
		},
		Action: func(_ context.Context, cmd *ufcli.Command) error {
			return run(cmd)
		},
	}
}

func run(cmd *ufcli.Command) error {
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
	target, err := store.GetByID(int64(id))
	if err != nil {
		if errors.Is(err, entry.ErrEntryNotFound) {
			output.Stdoutln("No entry deleted")
			return nil
		}
		return err
	}

	if target.DeletedAt != nil {
		output.Stdoutln("No entry deleted")
		return nil
	}

	shouldDelete, err := confirmDelete(cmd.Bool("yes"), target)
	if err != nil {
		return err
	}
	if !shouldDelete {
		output.Stdoutln("No entry deleted")
		return nil
	}

	deleted, err := store.SoftDeleteByID(int64(id))
	if err != nil {
		return err
	}

	if !deleted {
		output.Stdoutln("No entry deleted")
		return nil
	}

	output.Stdoutln(fmt.Sprintf("Deleted #%d", id))
	return nil
}

func confirmDelete(skipPrompt bool, target *entry.Entry) (bool, error) {
	if skipPrompt {
		return true, nil
	}
	if !output.IsInputTTY() {
		return false, errors.New("refusing to prompt on non-tty; pass --yes")
	}

	prompt := fmt.Sprintf(
		"Delete #%d (%s): %q?",
		target.ID,
		target.CreatedAt.Format("2006-01-02"),
		target.Text,
	)
	confirmed, err := output.Confirm(prompt)
	if err != nil {
		return false, err
	}

	return confirmed, nil
}
