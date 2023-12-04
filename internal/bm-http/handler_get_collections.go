package bmhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	bm "github.com/Tsapen/bm/internal/bm"
	"github.com/Tsapen/bm/pkg/api"
)

func parseGetCollectionsReq(r *http.Request) (*api.GetCollectionsReq, error) {
	q := r.URL.Query()
	req := &api.GetCollectionsReq{
		OrderBy: q.Get("order_by"),
	}

	var err error

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

func (b *serviceBundle) getCollections(ctx context.Context, r *api.GetCollectionsReq) (any, error) {
	collections, err := b.bookService.Collections(ctx, bm.CollectionsFilter(*r))
	if err != nil {
		return nil, fmt.Errorf("get collections: %w", err)
	}

	collectionsResp := make([]api.Collection, 0, len(collections))
	for _, c := range collections {
		collectionsResp = append(collectionsResp, api.Collection(c))
	}

	return &api.GetCollectionsResp{
		Collections: collectionsResp,
	}, nil
}
