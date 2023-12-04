package bmhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Tsapen/bm/pkg/api"
	"github.com/gorilla/mux"
)

func parseCreateBooksCollectionReq(r *http.Request) (*api.CreateBooksCollectionReq, error) {
	req := new(api.CreateBooksCollectionReq)
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

func (b *serviceBundle) createBooksCollection(ctx context.Context, r *api.CreateBooksCollectionReq) (any, error) {
	err := b.bookService.CreateBooksCollection(ctx, r.CID, r.BookIDs)
	if err != nil {
		return nil, fmt.Errorf("create books collection: %w", err)
	}

	return nil, nil
}
