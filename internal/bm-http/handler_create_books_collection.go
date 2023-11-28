package bmhttp

import (
	"context"
	"fmt"

	"github.com/Tsapen/bm/pkg/api"
)

func (b *serviceBundle) createBooksCollection(ctx context.Context, r *api.CreateBooksCollectionReq) (*api.CreateBooksCollectionResp, error) {
	err := b.bookService.CreateBooksCollection(ctx, r.CID, r.BookIDs)
	if err != nil {
		return nil, fmt.Errorf("create books collection: %w", err)
	}

	return &api.CreateBooksCollectionResp{
		Success: true,
	}, nil
}
