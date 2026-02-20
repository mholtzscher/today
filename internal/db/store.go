//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate -f ../../sqlc.yaml

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mholtzscher/today/internal/db/sqlc"
)

var ErrEntryNotFound = errors.New("entry not found")

// Entry is the domain model for a journal entry.
type Entry struct {
	ID         int64
	Text       string
	CreatedAt  time.Time
	ArchivedAt *time.Time
}

// Store wraps sqlc queries and provides domain-specific operations.
type Store struct {
	queries *sqlc.Queries
}

// NewStore creates a new Store from a database connection.
func NewStore(db *sql.DB) *Store {
	return &Store{
		queries: sqlc.New(db),
	}
}

// CreateEntry inserts a new entry with the given text.
func (s *Store) CreateEntry(ctx context.Context, text string) error {
	return s.queries.CreateEntry(ctx, text)
}

// CreateEntryAt inserts a new entry with an explicit creation time.
func (s *Store) CreateEntryAt(ctx context.Context, text string, createdAt time.Time) error {
	return s.queries.CreateEntryAt(ctx, sqlc.CreateEntryAtParams{
		Text:      text,
		CreatedAt: createdAt.Unix(),
	})
}

// GetEntry retrieves a single entry by ID.
func (s *Store) GetEntry(ctx context.Context, id int64) (*Entry, error) {
	row, err := s.queries.GetEntry(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEntryNotFound
		}
		return nil, fmt.Errorf("get entry: %w", err)
	}
	return convertEntry(row)
}

// ListEntries retrieves entries within the specified days window.
func (s *Store) ListEntries(ctx context.Context, days int, includeArchived bool) ([]Entry, error) {
	modifier := fmt.Sprintf("-%d days", days)

	var rows []sqlc.Entry
	var err error
	if includeArchived {
		rows, err = s.queries.ListEntriesSinceAll(ctx, modifier)
	} else {
		rows, err = s.queries.ListEntriesSince(ctx, modifier)
	}
	if err != nil {
		return nil, fmt.Errorf("list entries: %w", err)
	}

	return convertEntries(rows)
}

// ArchiveEntry marks an entry as archived.
func (s *Store) ArchiveEntry(ctx context.Context, id int64) (bool, error) {
	affected, err := s.queries.ArchiveEntry(ctx, id)
	if err != nil {
		return false, fmt.Errorf("archive entry: %w", err)
	}
	return affected == 1, nil
}

// RestoreEntry removes the archived status from an entry.
func (s *Store) RestoreEntry(ctx context.Context, id int64) (bool, error) {
	affected, err := s.queries.RestoreEntry(ctx, id)
	if err != nil {
		return false, fmt.Errorf("restore entry: %w", err)
	}
	return affected == 1, nil
}

// convertEntries converts a slice of sqlc entries to domain entries.
func convertEntries(rows []sqlc.Entry) ([]Entry, error) {
	entries := make([]Entry, len(rows))
	for i, row := range rows {
		entry, err := convertEntry(row)
		if err != nil {
			return nil, err
		}
		entries[i] = *entry
	}
	return entries, nil
}

// convertEntry converts a single sqlc entry to a domain entry.
func convertEntry(row sqlc.Entry) (*Entry, error) {
	e := &Entry{
		ID:        row.ID,
		Text:      row.Text,
		CreatedAt: time.Unix(row.CreatedAt, 0),
	}

	if row.ArchivedAt.Valid {
		archivedAt := time.Unix(row.ArchivedAt.Int64, 0)
		e.ArchivedAt = &archivedAt
	}

	return e, nil
}
