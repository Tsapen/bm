package bookservice

import (
	"context"
	"fmt"
	"time"

	bm "github.com/Tsapen/bm/internal/bm"
)

// Books retrieves a list of books based on the provided filter criteria.
func (s *Service) Books(ctx context.Context, f bm.BookFilter) ([]bm.Book, error) {
	switch f.OrderBy {
	case "id", "title", "author", "genre", "published_date", "edition":
	case "":
		f.OrderBy = "id"
	default:
		return nil, bm.ValidationError("incorrect order_by")
	}

	if f.Page < 0 {
		return nil, bm.ValidationError("incorrect page")
	}

	if f.Page == 0 {
		f.Page = 1
	}

	if f.PageSize < 0 {
		return nil, bm.ValidationError("page_size is negative")
	}

	if f.PageSize == 0 || f.PageSize > maxPageSize {
		f.PageSize = maxPageSize
	}

	books, err := s.db.Books(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("get books: %w", err)
	}

	return books, nil
}

// CreateBook creates a new book with the provided details.
func (s *Service) CreateBook(ctx context.Context, b bm.Book) (int64, error) {
	if b.Title == "" {
		return 0, bm.ValidationError("title is empty")
	}

	if b.Author == "" {
		return 0, bm.ValidationError("author is empty")
	}

	if b.Genre == "" {
		return 0, bm.ValidationError("genre is empty")
	}

	if !b.PublishedDate.IsZero() {
		b.PublishedDate = b.PublishedDate.Truncate(24 * time.Hour)
	}

	id, err := s.db.CreateBook(ctx, b)
	if err != nil {
		return 0, fmt.Errorf("create book: %w", err)
	}

	return id, nil
}

// UpdateBook updates an existing book with the provided details.
func (s *Service) UpdateBook(ctx context.Context, b bm.Book) error {
	if b.ID <= 0 {
		return bm.ValidationError("incorrect id")
	}

	if b.Title == "" {
		return bm.ValidationError("title is empty")
	}

	if b.Author == "" {
		return bm.ValidationError("author is empty")
	}

	if b.Genre == "" {
		return bm.ValidationError("genre is empty")
	}

	if !b.PublishedDate.IsZero() {
		b.PublishedDate = b.PublishedDate.Truncate(24 * time.Hour)
	}

	if err := s.db.UpdateBook(ctx, b); err != nil {
		return fmt.Errorf("update book: %w", err)
	}

	return nil
}

// DeleteBooks deletes multiple books based on their IDs.
func (s *Service) DeleteBooks(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return bm.ValidationError("ids list is empty")
	}

	if err := s.db.DeleteBooks(ctx, ids); err != nil {
		return fmt.Errorf("delete books: %w", err)
	}

	return nil
}
