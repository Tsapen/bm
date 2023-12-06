package bmhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	bm "github.com/Tsapen/bm/internal/bm"
	"github.com/Tsapen/bm/pkg/api"
	"github.com/gorilla/mux"
)

func parseUpdateBookReq(r *http.Request) (*api.UpdateBookReq, error) {
	req := new(api.UpdateBookReq)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}

	v := mux.Vars(r)

	reqID, err := strconv.ParseInt(v["book_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}

	req.ID = reqID

	return req, nil
}

func (b *serviceBundle) updateBook(ctx context.Context, r *api.UpdateBookReq) (any, error) {
	err := b.bookService.UpdateBook(ctx, bm.Book(*r))
	if err != nil {
		return nil, fmt.Errorf("update book: %w", err)
	}

	return nil, nil
}
