package bmhttp

import (
	"context"
	"fmt"
	"time"

	bm "github.com/Tsapen/bm/internal/bm"
)

type (
	updateBookReq struct {
		ID            int64     `json:"id"`
		Title         string    `json:"title"`
		Author        string    `json:"author"`
		PublishedDate time.Time `json:"published_date"`
		Edition       string    `json:"edition"`
		Description   string    `json:"description"`
		Genre         string    `json:"genre"`
	}

	updateBookResp struct {
		Success bool `json:"success"`
	}
)

func (b *serviceBundle) updateBook(ctx context.Context, r *updateBookReq) (any, error) {
	err := b.bookService.UpdateBook(ctx, bm.Book(*r))
	if err != nil {
		return nil, fmt.Errorf("update book: %w", err)
	}

	return updateBookResp{
		Success: true,
	}, nil
}
