package bookservice

import (
	"context"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
)

// Collection retrieves a collection by its id.
func (s *Service) Collection(ctx context.Context, cid int64) (*bm.Collection, error) {
	if cid < 0 {
		return nil, bm.NewValidationError("incorrect page")
	}

	collections, err := s.storage.Collection(ctx, cid)
	if err != nil {
		return nil, fmt.Errorf("get collections: %w", err)
	}

	return collections, nil
}

// Collections retrieves a list of collections based on the provided filter criteria.
func (s *Service) Collections(ctx context.Context, f bm.CollectionsFilter) ([]bm.Collection, error) {
	switch f.OrderBy {
	case "id", "name":
	case "":
		f.OrderBy = "id"
	default:
		return nil, bm.NewValidationError("incorrect order_by")
	}

	if f.Page < 0 {
		return nil, bm.NewValidationError("incorrect page")
	}

	if f.Page == 0 {
		f.Page = 1
	}

	if f.PageSize < 0 {
		return nil, bm.NewValidationError("page_size is negative")
	}

	if f.PageSize == 0 || f.PageSize > maxPageSize {
		f.PageSize = maxPageSize
	}

	collections, err := s.storage.Collections(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("get collections: %w", err)
	}

	return collections, nil
}

// CreateCollection creates a new collection with the provided details.
func (s *Service) CreateCollection(ctx context.Context, c bm.Collection) (int64, error) {
	if c.Name == "" {
		return 0, bm.NewValidationError("name is empty")
	}

	id, err := s.storage.CreateCollection(ctx, c)
	if err != nil {
		return 0, fmt.Errorf("create collection: %w", err)
	}

	return id, nil
}

// CreateBooksCollection adds a list of books to an existing collection.
func (s *Service) CreateBooksCollection(ctx context.Context, cID int64, bookIDs []int64) error {
	if cID <= 0 {
		return bm.NewValidationError("incorrect collection_id")
	}

	if len(bookIDs) == 0 {
		return bm.NewValidationError("empty book ids list")
	}

	if err := s.storage.CreateBooksCollection(ctx, cID, bookIDs); err != nil {
		return fmt.Errorf("add books to collection: %w", err)
	}

	return nil
}

// DeleteBooksCollection removes a list of books from an existing collection.
func (s *Service) DeleteBooksCollection(ctx context.Context, cID int64, bookIDs []int64) error {
	if cID <= 0 {
		return bm.NewValidationError("incorrect collection_id")
	}

	if len(bookIDs) == 0 {
		return bm.NewValidationError("empty book ids list")
	}

	if err := s.storage.DeleteBooksCollection(ctx, cID, bookIDs); err != nil {
		return fmt.Errorf("remove books from collection: %w", err)
	}

	return nil
}

// UpdateCollection updates an existing collection with the provided details.
func (s *Service) UpdateCollection(ctx context.Context, c bm.Collection) error {
	if c.ID <= 0 {
		return bm.NewValidationError("incorrect id")
	}

	if c.Name == "" {
		return bm.NewValidationError("name is empty")
	}

	if err := s.storage.UpdateCollection(ctx, c); err != nil {
		return fmt.Errorf("update collection: %w", err)
	}

	return nil
}

// DeleteCollection deletes a collection based on its ID.
func (s *Service) DeleteCollection(ctx context.Context, cID int64) error {
	if cID <= 0 {
		return bm.NewValidationError("incorrect id")
	}

	if err := s.storage.DeleteCollection(ctx, cID); err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}

	return nil
}
