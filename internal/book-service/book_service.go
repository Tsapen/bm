package bookservice

import (
	bm "github.com/Tsapen/bm/internal/bm"
)

const (
	maxPageSize = 50
)

// Service stores and manages books and collections.
type Service struct {
	storage bm.Storage
}

// New constructs new book service.
func New(db bm.Storage) *Service {
	return &Service{
		storage: db,
	}
}
