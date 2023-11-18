package bmhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

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

func parseGetBooksReq(r *http.Request) (*getBooksReq, error) {
	vars := mux.Vars(r)
	var startDate, finishDate time.Time
	var desc bool
	var page, pageSize int64

	cID, err := strconv.ParseInt(vars["collection_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("incorrect collection_id: %w", err)
	}

	if startDateStr, ok := vars["start_date"]; ok {
		startDate, err = time.Parse(time.DateOnly, startDateStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect start_date: %w", err)
		}
	}

	if finishDateStr, ok := vars["finish_date"]; ok {
		finishDate, err = time.Parse(time.DateOnly, finishDateStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect finish_date: %w", err)
		}
	}

	if descStr, ok := vars["desc"]; ok {
		desc, err = strconv.ParseBool(descStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect desc: %w", err)
		}
	}

	if pageStr, ok := vars["page"]; ok {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect page: %w", err)
		}
	}

	if pageSizeStr, ok := vars["page_size"]; ok {
		page, err = strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect page_size: %w", err)
		}
	}

	return &getBooksReq{
		Author:       vars["author"],
		Genre:        vars["genre"],
		OrderBy:      vars["order_by"],
		Desc:         desc,
		StartDate:    startDate,
		FinishDate:   finishDate,
		CollectionID: cID,
		Page:         page,
		PageSize:     pageSize,
	}, nil
}

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
