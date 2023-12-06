package bmhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	bm "github.com/Tsapen/bm/internal/bm"
	"github.com/Tsapen/bm/pkg/api"
)

func parseGetBooksReq(r *http.Request) (*api.GetBooksReq, error) {
	q := r.URL.Query()
	req := &api.GetBooksReq{
		Author:  q.Get("author"),
		Genre:   q.Get("genre"),
		OrderBy: q.Get("order_by"),
	}

	var err error

	if cidStr := q.Get("collection_id"); cidStr != "" {
		req.CollectionID, err = strconv.ParseInt(cidStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect collection_id: %w", err)
		}
	}

	if startDateStr := q.Get("start_date"); startDateStr != "" {
		req.StartDate, err = time.Parse(time.DateOnly, startDateStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect start_date: %w", err)
		}
	}

	if finishDateStr := q.Get("finish_date"); finishDateStr != "" {
		req.FinishDate, err = time.Parse(time.DateOnly, finishDateStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect finish_date: %w", err)
		}
	}

	if descStr := q.Get("desc"); descStr != "" {
		req.Desc, err = strconv.ParseBool(descStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect desc: %w", err)
		}
	}

	if pageStr := q.Get("page"); pageStr != "" {
		req.Page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect page: %w", err)
		}
	}

	if pageSizeStr := q.Get("page_size"); pageSizeStr != "" {
		req.PageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect page_size: %w", err)
		}
	}

	return req, nil
}

func (b *serviceBundle) getBooks(ctx context.Context, r *api.GetBooksReq) (any, error) {
	f := bm.BookFilter{
		Author:       r.Author,
		Genre:        r.Genre,
		CollectionID: r.CollectionID,
		StartDate:    r.StartDate,
		FinishDate:   r.FinishDate,
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
