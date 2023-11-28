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
	var startDate, finishDate time.Time
	var desc bool
	var id, cID, page, pageSize int64
	var err error

	q := r.URL.Query()
	if idStr := q.Get("id"); idStr != "" {
		id, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect id: %w", err)
		}
	}

	if cidStr := q.Get("collection_id"); cidStr != "" {
		cID, err = strconv.ParseInt(cidStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect collection_id: %w", err)
		}
	}

	if startDateStr := q.Get("start_date"); startDateStr != "" {
		startDate, err = time.Parse(time.DateOnly, startDateStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect start_date: %w", err)
		}
	}

	if finishDateStr := q.Get("finish_date"); finishDateStr != "" {
		finishDate, err = time.Parse(time.DateOnly, finishDateStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect finish_date: %w", err)
		}
	}

	if descStr := q.Get("desc"); descStr != "" {
		desc, err = strconv.ParseBool(descStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect desc: %w", err)
		}
	}

	if pageStr := q.Get("page"); pageStr != "" {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect page: %w", err)
		}
	}

	if pageSizeStr := q.Get("page_size"); pageSizeStr != "" {
		pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect page_size: %w", err)
		}
	}

	return &api.GetBooksReq{
		ID:           id,
		Author:       q.Get("author"),
		Genre:        q.Get("genre"),
		StartDate:    startDate,
		FinishDate:   finishDate,
		CollectionID: cID,
		OrderBy:      q.Get("order_by"),
		Desc:         desc,
		Page:         page,
		PageSize:     pageSize,
	}, nil
}

func (b *serviceBundle) getBooks(ctx context.Context, r *api.GetBooksReq) (*api.GetBooksResp, error) {
	f := bm.BookFilter{
		ID:           r.ID,
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
