package bmhttp

import (
	"context"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
)

type (
	updateCollectionReq struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"decription"`
	}

	updateCollectionResp struct {
		Success bool `json:"success"`
	}
)

func (b *serviceBundle) updateCollection(ctx context.Context, r *updateCollectionReq) (any, error) {
	err := b.bookService.UpdateCollection(ctx, bm.Collection(*r))
	if err != nil {
		return nil, fmt.Errorf("update collection: %w", err)
	}

	return updateCollectionResp{
		Success: true,
	}, nil
}
