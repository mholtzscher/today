package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite" // sqlite driver for database/sql

	"github.com/mholtzscher/today/internal/db/migrations"
	"github.com/mholtzscher/today/internal/output"
)

type silentLogger struct{}

func (l *silentLogger) Printf(_ string, _ ...any) {}

func (l *silentLogger) Fatalf(format string, v ...any) {
	output.Stderrln(fmt.Sprintf(format, v...))
	os.Exit(1)
}

var _ goose.Logger = (*silentLogger)(nil)

func Open(path string) (*sql.DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if migrateErr := migrate(db); migrateErr != nil {
		_ = db.Close()
		return nil, fmt.Errorf("run migrations: %w", migrateErr)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(&silentLogger{})
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

var _ = io.EOF
