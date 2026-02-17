package entry

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrEntryNotFound = errors.New("entry not found")

type Entry struct {
	ID        int64
	Text      string
	CreatedAt time.Time
	DeletedAt *time.Time
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

func (s *Store) GetByDays(days int, includeDeleted bool) ([]Entry, error) {
	query := `
		SELECT id, text, created_at, deleted_at
		FROM entries
		WHERE date(created_at) >= date('now', '-' || ? || ' days')
	`
	if !includeDeleted {
		query += "\n\t\tAND deleted_at IS NULL"
	}
	query += "\n\t\tORDER BY created_at DESC"

	rows, err := s.db.QueryContext(context.Background(), query, days)
	if err != nil {
		return nil, fmt.Errorf("query entries: %w", err)
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		var createdAtStr string
		var deletedAt sql.NullString
		if scanErr := rows.Scan(&e.ID, &e.Text, &createdAtStr, &deletedAt); scanErr != nil {
			return nil, fmt.Errorf("scan entry: %w", scanErr)
		}
		e.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("parse created_at: %w", err)
		}
		if deletedAt.Valid {
			parsedDeletedAt, parseErr := time.Parse("2006-01-02 15:04:05", deletedAt.String)
			if parseErr != nil {
				return nil, fmt.Errorf("parse deleted_at: %w", parseErr)
			}
			e.DeletedAt = &parsedDeletedAt
		}
		entries = append(entries, e)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("rows error: %w", rowsErr)
	}
	return entries, nil
}

func (s *Store) GetByID(id int64) (*Entry, error) {
	query := `
		SELECT id, text, created_at, deleted_at
		FROM entries
		WHERE id = ?
		LIMIT 1
	`
	row := s.db.QueryRowContext(context.Background(), query, id)

	var e Entry
	var createdAtStr string
	var deletedAt sql.NullString
	if err := row.Scan(&e.ID, &e.Text, &createdAtStr, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrEntryNotFound
		}
		return nil, fmt.Errorf("query entry by id: %w", err)
	}

	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("parse created_at: %w", err)
	}
	e.CreatedAt = createdAt

	if deletedAt.Valid {
		parsedDeletedAt, parseErr := time.Parse("2006-01-02 15:04:05", deletedAt.String)
		if parseErr != nil {
			return nil, fmt.Errorf("parse deleted_at: %w", parseErr)
		}
		e.DeletedAt = &parsedDeletedAt
	}

	return &e, nil
}

func (s *Store) SoftDeleteByID(id int64) (bool, error) {
	result, err := s.db.ExecContext(
		context.Background(),
		"UPDATE entries SET deleted_at = datetime('now') WHERE id = ? AND deleted_at IS NULL",
		id,
	)
	if err != nil {
		return false, fmt.Errorf("soft delete entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("soft delete rows affected: %w", err)
	}

	return rowsAffected == 1, nil
}

func (s *Store) RestoreByID(id int64) (bool, error) {
	result, err := s.db.ExecContext(
		context.Background(),
		"UPDATE entries SET deleted_at = NULL WHERE id = ? AND deleted_at IS NOT NULL",
		id,
	)
	if err != nil {
		return false, fmt.Errorf("restore entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("restore rows affected: %w", err)
	}

	return rowsAffected == 1, nil
}
