package bmhttp

import (
	"context"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
	"github.com/Tsapen/bm/pkg/api"
)

func (b *serviceBundle) createBook(ctx context.Context, r *api.CreateBookReq) (any, error) {
	bookData := bm.Book{
		Title:         r.Title,
		Author:        r.Author,
		PublishedDate: r.PublishedDate,
		Edition:       r.Edition,
		Description:   r.Description,
		Genre:         r.Genre,
	}

	id, err := b.bookService.CreateBook(ctx, bookData)
	if err != nil {
		return nil, fmt.Errorf("create book: %w", err)
	}

	return &api.CreateBookResp{
		ID: id,
	}, nil
}
