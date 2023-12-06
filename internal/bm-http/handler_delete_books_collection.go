package bmhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Tsapen/bm/pkg/api"
)

func parseDeleteBooksCollectionReq(r *http.Request) (*api.DeleteBooksCollectionReq, error) {
	req := new(api.DeleteBooksCollectionReq)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}

	v := mux.Vars(r)

	reqID, err := strconv.ParseInt(v["collection_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}

	req.CID = reqID

	return req, nil
}

func (b *serviceBundle) deleteBooksCollection(ctx context.Context, r *api.DeleteBooksCollectionReq) (any, error) {
	err := b.bookService.DeleteBooksCollection(ctx, r.CID, r.BookIDs)
	if err != nil {
		return nil, fmt.Errorf("delete books collection: %w", err)
	}

	return nil, nil
}
