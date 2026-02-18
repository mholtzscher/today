package entry

import (
	"context"
	"errors"

	"github.com/mholtzscher/today/internal/db"
)

// Service defines the interface for entry operations.
type Service interface {
	CreateEntry(ctx context.Context, text string) error
	GetEntry(ctx context.Context, id int64) (*db.Entry, error)
	ListEntries(ctx context.Context, days int, includeArchived bool) ([]db.Entry, error)
	ArchiveEntry(ctx context.Context, id int64) (bool, error)
	RestoreEntry(ctx context.Context, id int64) (bool, error)
}

// Ensure Service implements the interface.
var _ Service = (*service)(nil)

// service provides business logic for entry operations.
type service struct {
	store *db.Store
}

// NewService creates a new entry service.
func NewService(store *db.Store) Service {
	return &service{store: store}
}

// CreateEntry adds a new journal entry.
func (s *service) CreateEntry(ctx context.Context, text string) error {
	return s.store.CreateEntry(ctx, text)
}

// GetEntry retrieves a single entry by ID.
func (s *service) GetEntry(ctx context.Context, id int64) (*db.Entry, error) {
	return s.store.GetEntry(ctx, id)
}

// ListEntries retrieves entries for the specified time period.
func (s *service) ListEntries(ctx context.Context, days int, includeArchived bool) ([]db.Entry, error) {
	return s.store.ListEntries(ctx, days, includeArchived)
}

// ArchiveEntry archives an entry and returns true if it was archived.
func (s *service) ArchiveEntry(ctx context.Context, id int64) (bool, error) {
	entry, err := s.store.GetEntry(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}
	if entry.ArchivedAt != nil {
		return false, nil
	}
	return s.store.ArchiveEntry(ctx, id)
}

// RestoreEntry restores an archived entry and returns true if it was restored.
func (s *service) RestoreEntry(ctx context.Context, id int64) (bool, error) {
	entry, err := s.store.GetEntry(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}
	if entry.ArchivedAt == nil {
		return false, nil
	}
	return s.store.RestoreEntry(ctx, id)
}
