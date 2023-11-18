package bookservice

import (
	"context"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
)

func (s *Service) Collections(ctx context.Context, cf bm.CollectionsFilter) ([]bm.Collection, error) {
	collections, err := s.db.Collections(ctx, cf)
	if err != nil {
		return nil, fmt.Errorf("get collections: %w", err)
	}

	return collections, nil
}

func (s *Service) GetCollectionInfo(ctx context.Context, f bm.BooksCollectionFilter) (bm.CollectionInfo, error) {
	cis, err := s.db.Collections(ctx, bm.CollectionsFilter{IDs: []int64{f.CID}})
	if err != nil {
		return bm.CollectionInfo{}, fmt.Errorf("get collection: %w", err)
	}

	if len(cis) != 1 {
		return bm.CollectionInfo{}, fmt.Errorf("collection not found")
	}

	ci := bm.CollectionInfo{
		Collection: cis[0],
	}

	bf := bm.BookFilter{
		CollectionID: f.CID,
		Desc:         f.Desc,
		OrderBy:      f.OrderBy,
		Page:         f.Page,
		PageSize:     f.PageSize,
	}
	ci.Books, err = s.db.Books(ctx, bf)
	if err != nil {
		return bm.CollectionInfo{}, fmt.Errorf("get books: %w", err)
	}

	return ci, nil
}

func (s *Service) CreateCollection(ctx context.Context, c bm.Collection) (int64, error) {
	id, err := s.db.CreateCollection(ctx, c)
	if err != nil {
		return 0, fmt.Errorf("create collection: %w", err)
	}

	return id, nil
}

func (s *Service) CreateBooksCollection(ctx context.Context, cID int64, bookIDs []int64) error {
	if err := s.db.CreateBooksCollection(ctx, cID, bookIDs); err != nil {
		return fmt.Errorf("add books to collection: %w", err)
	}

	return nil
}

func (s *Service) DeleteBooksCollection(ctx context.Context, cID int64, bookIDs []int64) error {
	if err := s.db.DeleteBooksCollection(ctx, cID, bookIDs); err != nil {
		return fmt.Errorf("remove books from collection: %w", err)
	}

	return nil
}

func (s *Service) UpdateCollection(ctx context.Context, collection bm.Collection) error {
	if err := s.db.UpdateCollection(ctx, collection); err != nil {
		return fmt.Errorf("update collection: %w", err)
	}

	return nil
}

func (s *Service) DeleteCollection(ctx context.Context, cID int64) error {
	if err := s.db.DeleteCollection(ctx, cID); err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}

	return nil
}
