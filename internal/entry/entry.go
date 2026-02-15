package entry

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Entry struct {
	ID        int64
	Text      string
	CreatedAt time.Time
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Insert(text string) error {
	_, err := s.db.ExecContext(
		context.Background(),
		"INSERT INTO entries (text) VALUES (?)",
		text,
	)
	if err != nil {
		return fmt.Errorf("insert entry: %w", err)
	}
	return nil
}

func (s *Store) GetByDays(days int) ([]Entry, error) {
	query := `
		SELECT id, text, created_at
		FROM entries
		WHERE date(created_at) >= date('now', '-' || ? || ' days')
		ORDER BY created_at DESC
	`
	rows, err := s.db.QueryContext(context.Background(), query, days)
	if err != nil {
		return nil, fmt.Errorf("query entries: %w", err)
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		var createdAtStr string
		if scanErr := rows.Scan(&e.ID, &e.Text, &createdAtStr); scanErr != nil {
			return nil, fmt.Errorf("scan entry: %w", scanErr)
		}
		e.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("parse created_at: %w", err)
		}
		entries = append(entries, e)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("rows error: %w", rowsErr)
	}
	return entries, nil
}
