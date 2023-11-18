package bookservice

import (
	bm "github.com/Tsapen/bm/internal/bm"
)

type Service struct {
	db bm.Storage
}

func New(db bm.Storage) *Service {
	return &Service{
		db: db,
	}
}
