package bmunixsocket

import (
	"context"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
	"github.com/Tsapen/bm/pkg/api"
)

func (b *serviceBundle) getBooks(ctx context.Context, r *api.GetBooksReq) (*api.GetBooksResp, error) {
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

	booksResp := make([]api.Book, 0, len(books))
	for _, b := range books {
		booksResp = append(booksResp, api.Book(b))
	}

	return &api.GetBooksResp{
		Books: booksResp,
	}, nil
}
