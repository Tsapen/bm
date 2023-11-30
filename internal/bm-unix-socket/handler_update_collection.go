package bmunixsocket

import (
	"context"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
	"github.com/Tsapen/bm/pkg/api"
)

func (b *serviceBundle) updateCollection(ctx context.Context, r *api.UpdateCollectionReq) (*api.UpdateCollectionResp, error) {
	err := b.bookService.UpdateCollection(ctx, bm.Collection(*r))
	if err != nil {
		return nil, fmt.Errorf("update collection: %w", err)
	}

	return &api.UpdateCollectionResp{
		Success: true,
	}, nil
}
