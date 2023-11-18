package bmunixsocket

import (
	"context"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
)

type (
	createCollectionReq struct {
		Name        string `json:"name"`
		Description string `json:"decription"`
	}

	createCollectionResp struct {
		ID int64 `json:"id"`
	}
)

func (b *serviceBundle) createCollection(ctx context.Context, r *createCollectionReq) (any, error) {
	collectionData := bm.Collection{
		Name:        r.Name,
		Description: r.Description,
	}

	id, err := b.bookService.CreateCollection(ctx, collectionData)
	if err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}

	return createCollectionResp{
		ID: id,
	}, nil
}
