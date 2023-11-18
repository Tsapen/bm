package bmhttp

import (
	"context"
	"fmt"
)

type (
	deleteBooksReq struct {
		IDs []int64 `json:"ids"`
	}

	deleteBooksResp struct {
		Success bool `json:"success"`
	}
)

func (b *serviceBundle) deleteBooks(ctx context.Context, r *deleteBooksReq) (any, error) {
	err := b.bookService.DeleteBooks(ctx, r.IDs)
	if err != nil {
		return nil, fmt.Errorf("delete book: %w", err)
	}

	return deleteBooksResp{
		Success: true,
	}, nil
}
