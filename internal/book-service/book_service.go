package bookservice

import (
	bm "github.com/Tsapen/bm/internal/bm"
)

const (
	maxPageSize = 50
)

type Service struct {
	db bm.Storage
}

func New(db bm.Storage) *Service {
	return &Service{
		db: db,
	}
}
