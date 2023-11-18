package bmunixsocket

import (
	"context"
	"fmt"
	"time"

	bm "github.com/Tsapen/bm/internal/bm"
)

type (
	getBooksReq struct {
		Author       string
		Genre        string
		CollectionID int64
		StartDate    time.Time
		FinishDate   time.Time
		OrderBy      string
		Desc         bool
		Page         int64
		PageSize     int64
	}

	book struct {
		ID            int64     `json:"id"`
		Title         string    `json:"title"`
		Author        string    `json:"author"`
		PublishedDate time.Time `json:"published_date"`
		Edition       string    `json:"edition"`
		Description   string    `json:"description"`
		Genre         string    `json:"genre"`
	}

	getBookResp struct {
		Books []book `json:"books"`
	}
)

func (b *serviceBundle) getBooks(ctx context.Context, r *getBooksReq) (any, error) {
	f := bm.BookFilter{
		Author:       r.Author,
		Genre:        r.Genre,
		CollectionID: r.CollectionID,
		OrderBy:      r.OrderBy,
		Desc:         r.Desc,
		Page:         r.Page,
		PageSize:     r.PageSize,
	}

	books, err := b.bookService.Books(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("get books: %w", err)
	}

	booksResp := make([]book, 0, len(books))
	for _, b := range books {
		booksResp = append(booksResp, book(b))
	}

	return getBookResp{
		Books: booksResp,
	}, nil
}
