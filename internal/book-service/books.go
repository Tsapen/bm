package bookservice

import (
	"context"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
)

func (s *Service) Books(ctx context.Context, bp bm.BookFilter) ([]bm.Book, error) {
	books, err := s.db.Books(ctx, bp)
	if err != nil {
		return nil, fmt.Errorf("get books: %w", err)
	}

	return books, nil
}

func (s *Service) CreateBook(ctx context.Context, book bm.Book) (int64, error) {
	ids, err := s.db.CreateBook(ctx, book)
	if err != nil {
		return 0, fmt.Errorf("create book: %w", err)
	}

	return ids, nil
}

func (s *Service) UpdateBook(ctx context.Context, book bm.Book) error {
	if err := s.db.UpdateBook(ctx, book); err != nil {
		return fmt.Errorf("update book: %w", err)
	}

	return nil
}

func (s *Service) DeleteBooks(ctx context.Context, ids []int64) error {
	if err := s.db.DeleteBooks(ctx, ids); err != nil {
		return fmt.Errorf("delete books: %w", err)
	}

	return nil
}
