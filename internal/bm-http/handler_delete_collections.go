package bmhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Tsapen/bm/pkg/api"
)

func parseDeleteCollectionReq(r *http.Request) (*api.DeleteCollectionReq, error) {
	v := mux.Vars(r)

	req := new(api.DeleteCollectionReq)
	var err error
	req.ID, err = strconv.ParseInt(v["collection_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}

	return req, nil
}

func (b *serviceBundle) deleteCollection(ctx context.Context, r *api.DeleteCollectionReq) (any, error) {
	err := b.bookService.DeleteCollection(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("delete collection: %w", err)
	}

	return nil, nil
}
