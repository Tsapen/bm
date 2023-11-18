package bmhttp

import (
	"context"
	"fmt"
)

type (
	deleteCollectionReq struct {
		ID int64 `json:"id"`
	}

	deleteCollectionResp struct {
		Success bool `json:"success"`
	}
)

func (b *serviceBundle) deleteCollection(ctx context.Context, r *deleteCollectionReq) (any, error) {
	err := b.bookService.DeleteCollection(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("delete collection: %w", err)
	}

	return deleteCollectionResp{
		Success: true,
	}, nil
}
