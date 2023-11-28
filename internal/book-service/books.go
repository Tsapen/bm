package bookservice

import (
	"context"
	"fmt"
	"time"

	bm "github.com/Tsapen/bm/internal/bm"
)

func (s *Service) Books(ctx context.Context, f bm.BookFilter) ([]bm.Book, error) {
	if f.OrderBy == "" {
		f.OrderBy = "id"
	}

	if f.Page == 0 {
		f.Page = 1
	}

	if f.PageSize == 0 {
		f.PageSize = 50
	}

	books, err := s.db.Books(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("get books: %w", err)
	}

	return books, nil
}

func (s *Service) CreateBook(ctx context.Context, book bm.Book) (int64, error) {
	if !book.PublishedDate.IsZero() {
		book.PublishedDate = book.PublishedDate.Truncate(24 * time.Hour)
	}

	ids, err := s.db.CreateBook(ctx, book)
	if err != nil {
		return 0, fmt.Errorf("create book: %w", err)
	}

	return ids, nil
}

func (s *Service) UpdateBook(ctx context.Context, book bm.Book) error {
	if !book.PublishedDate.IsZero() {
		book.PublishedDate = book.PublishedDate.Truncate(24 * time.Hour)
	}

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
