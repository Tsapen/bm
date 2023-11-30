package bmunixsocket

import (
	"context"
	"fmt"

	"github.com/Tsapen/bm/pkg/api"
)

func (b *serviceBundle) deleteCollection(ctx context.Context, r *api.DeleteCollectionReq) (*api.DeleteCollectionResp, error) {
	err := b.bookService.DeleteCollection(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("delete collection: %w", err)
	}

	return &api.DeleteCollectionResp{
		Success: true,
	}, nil
}
