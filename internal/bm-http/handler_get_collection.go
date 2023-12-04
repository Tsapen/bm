package bmhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Tsapen/bm/pkg/api"
)

func parseGetCollectionReq(r *http.Request) (*api.GetCollectionReq, error) {
	req := &api.GetCollectionReq{}

	var err error
	v := mux.Vars(r)
	if idStr := v["collection_id"]; idStr != "" {
		req.ID, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect id: %w", err)
		}
	}

	return req, nil
}

func (b *serviceBundle) getCollection(ctx context.Context, r *api.GetCollectionReq) (any, error) {
	collection, err := b.bookService.Collection(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}

	return &api.GetCollectionResp{
		Collection: api.Collection(*collection),
	}, nil
}
