package bmhttp

import (
	"context"
	"fmt"

	"github.com/Tsapen/bm/pkg/api"
)

func (b *serviceBundle) deleteBooks(ctx context.Context, r *api.DeleteBooksReq) (any, error) {
	err := b.bookService.DeleteBooks(ctx, r.IDs)
	if err != nil {
		return nil, fmt.Errorf("delete book: %w", err)
	}

	return nil, nil
}
