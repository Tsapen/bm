package bmhttp

import (
	"context"
	"fmt"

	"github.com/Tsapen/bm/pkg/api"
)

func (b *serviceBundle) deleteBooksCollection(ctx context.Context, r *api.DeleteBooksCollectionReq) (*api.DeleteBooksCollectionResp, error) {
	err := b.bookService.DeleteBooksCollection(ctx, r.CID, r.BookIDs)
	if err != nil {
		return nil, fmt.Errorf("delete books collection: %w", err)
	}

	return &api.DeleteBooksCollectionResp{
		Success: true,
	}, nil
}
