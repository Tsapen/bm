package bmunixsocket

import (
	"context"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
	"github.com/Tsapen/bm/pkg/api"
)

func (b *serviceBundle) updateBook(ctx context.Context, r *api.UpdateBookReq) (*api.UpdateBookResp, error) {
	err := b.bookService.UpdateBook(ctx, bm.Book(*r))
	if err != nil {
		return nil, fmt.Errorf("update book: %w", err)
	}

	return &api.UpdateBookResp{
		Success: true,
	}, nil
}
