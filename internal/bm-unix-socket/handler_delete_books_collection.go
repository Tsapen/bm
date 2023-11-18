package bmunixsocket

import (
	"context"
	"fmt"
)

type (
	deleteBooksCollectionReq struct {
		CID     int64   `json:"collection_id"`
		BookIDs []int64 `json:"books_ids"`
	}

	deleteBooksCollectionResp struct {
		Success bool `json:"success"`
	}
)

func (b *serviceBundle) deleteBooksCollection(ctx context.Context, r *deleteBooksCollectionReq) (any, error) {
	err := b.bookService.CreateBooksCollection(ctx, r.CID, r.BookIDs)
	if err != nil {
		return nil, fmt.Errorf("delete books collection: %w", err)
	}

	return &deleteBooksCollectionResp{
		Success: true,
	}, nil
}
