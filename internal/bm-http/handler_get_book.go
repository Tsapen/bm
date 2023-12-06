package bmhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Tsapen/bm/pkg/api"
)

func parseGetBookReq(r *http.Request) (*api.GetBookReq, error) {
	req := &api.GetBookReq{}

	var err error
	v := mux.Vars(r)
	if idStr := v["book_id"]; idStr != "" {
		req.ID, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect id: %w", err)
		}
	}

	return req, nil
}

func (b *serviceBundle) getBook(ctx context.Context, r *api.GetBookReq) (any, error) {
	book, err := b.bookService.Book(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("get book: %w", err)
	}

	return &api.GetBookResp{
		Book: api.Book(*book),
	}, nil
}
