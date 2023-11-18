package bmunixsocket

import (
	"context"
	"fmt"
)

type (
	createBooksCollectionReq struct {
		CID     int64   `json:"collection_id"`
		BookIDs []int64 `json:"books_ids"`
	}

	createBooksCollectionResp struct {
		Success bool `json:"success"`
	}
)

func (b *serviceBundle) createBooksCollection(ctx context.Context, r *createBooksCollectionReq) (any, error) {
	err := b.bookService.CreateBooksCollection(ctx, r.CID, r.BookIDs)
	if err != nil {
		return nil, fmt.Errorf("create books collection: %w", err)
	}

	return &createBooksCollectionResp{
		Success: true,
	}, nil
}
