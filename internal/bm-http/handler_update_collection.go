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

func parseUpdateCollectionReq(r *http.Request) (*api.UpdateCollectionReq, error) {
	req := new(api.UpdateCollectionReq)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}

	v := mux.Vars(r)

	reqID, err := strconv.ParseInt(v["collection_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}

	req.ID = reqID

	return req, nil
}

func (b *serviceBundle) updateCollection(ctx context.Context, r *api.UpdateCollectionReq) (any, error) {
	err := b.bookService.UpdateCollection(ctx, bm.Collection(*r))
	if err != nil {
		return nil, fmt.Errorf("update collection: %w", err)
	}

	return nil, nil
}
