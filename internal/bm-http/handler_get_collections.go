package bmhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	bm "github.com/Tsapen/bm/internal/bm"
)

type (
	getCollectionsReq struct {
		IDs      []int64
		OrderBy  string
		Desc     bool
		Page     int64
		PageSize int64
	}

	collection struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"decription"`
	}

	getCollectionsResp struct {
		Collections []collection `json:"collections"`
	}
)

func parseGetCollectionsReq(r *http.Request) (*getCollectionsReq, error) {
	vars := mux.Vars(r)

	var desc bool
	var page, pageSize int64
	var err error

	idStrs := strings.Split(vars["ids"], ",")
	ids := make([]int64, 0, len(idStrs))
	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect id: %w", err)
		}

		ids = append(ids, id)
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

	return &getCollectionsReq{
		IDs:      ids,
		OrderBy:  vars["order_by"],
		Desc:     desc,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (b *serviceBundle) getCollections(ctx context.Context, r *getCollectionsReq) (any, error) {
	collections, err := b.bookService.Collections(ctx, bm.CollectionsFilter(*r))
	if err != nil {
		return nil, fmt.Errorf("get collections: %w", err)
	}

	collectionsResp := make([]collection, 0, len(collections))
	for _, c := range collections {
		collectionsResp = append(collectionsResp, collection(c))
	}

	return getCollectionsResp{
		Collections: collectionsResp,
	}, nil
}
