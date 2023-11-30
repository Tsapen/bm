package bmunixsocket

import (
	"context"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
	"github.com/Tsapen/bm/pkg/api"
)

func (b *serviceBundle) createCollection(ctx context.Context, r *api.CreateCollectionReq) (*api.CreateCollectionResp, error) {
	collectionData := bm.Collection{
		Name:        r.Name,
		Description: r.Description,
	}

	id, err := b.bookService.CreateCollection(ctx, collectionData)
	if err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}

	return &api.CreateCollectionResp{
		ID: id,
	}, nil
}
