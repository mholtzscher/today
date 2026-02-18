package archive

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
		Name:  "archive",
		Usage: "Archive an entry by id",
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
			output.Stdoutln("No entry archived")
			return nil
		}
		return err
	}

	if target.ArchivedAt != nil {
		output.Stdoutln("No entry archived")
		return nil
	}

	shouldArchive, err := confirmArchive(cmd.Bool("yes"), target)
	if err != nil {
		return err
	}
	if !shouldArchive {
		output.Stdoutln("No entry archived")
		return nil
	}

	archived, err := store.ArchiveByID(int64(id))
	if err != nil {
		return err
	}

	if !archived {
		output.Stdoutln("No entry archived")
		return nil
	}

	output.Stdoutln(fmt.Sprintf("Archived #%d", id))
	return nil
}

func confirmArchive(skipPrompt bool, target *entry.Entry) (bool, error) {
	if skipPrompt {
		return true, nil
	}
	if !output.IsInputTTY() {
		return false, errors.New("refusing to prompt on non-tty; pass --yes")
	}

	prompt := fmt.Sprintf(
		"Archive #%d (%s): %q?",
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
