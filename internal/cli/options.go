package cli

import (
	"context"
	"os"
	"path/filepath"
)

const (
	FlagVerbose = "verbose"
	FlagNoColor = "no-color"
	FlagDB      = "db"
)

type GlobalOptions struct {
	Verbose bool
	NoColor bool
}

func GlobalOptionsFromContext(_ context.Context) GlobalOptions {
	return GlobalOptions{
		Verbose: false,
		NoColor: false,
	}
}

func DefaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "today.db"
	}
	return filepath.Join(home, "today.db")
}

func GetDBPath(cmd interface{ String(name string) string }) string {
	path := cmd.String(FlagDB)
	if path == "" {
		return DefaultDBPath()
	}
	return path
}
