package bmunixsocket

import (
	"context"
	"time"

	bm "github.com/Tsapen/bm/internal/bm"
)

type (
	createBookReq struct {
		Title         string    `json:"title"`
		Author        string    `json:"author"`
		PublishedDate time.Time `json:"published_date"`
		Edition       string    `json:"edition"`
		Description   string    `json:"description"`
		Genre         string    `json:"genre"`
	}

	createBookResp struct {
		ID int64 `json:"id"`
	}
)

func (b *serviceBundle) createBook(ctx context.Context, r *createBookReq) (any, error) {
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
		return nil, err
	}

	return createBookResp{
		ID: id,
	}, nil
}
