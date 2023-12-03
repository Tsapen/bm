package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	bmconfig "github.com/Tsapen/bm/internal/config"
	"github.com/Tsapen/bm/pkg/api"
	httpclient "github.com/Tsapen/bm/pkg/http-client"
)

type (
	flags struct {
		Action        *string
		ID            *int64
		IDs           *string
		Title         *string
		Author        *string
		PublishedDate *string
		Edition       *string
		Description   *string
		Genre         *string
		CollectionID  *int64
		BookIDs       *string
		Name          *string
		StartDate     *string
		FinishedDate  *string
		OrderBy       *string
		Desc          *bool
		Page          *int64
		PageSize      *int64
	}

	config struct {
		SocketPath string
		Timeout    time.Duration
	}

	cliClient struct {
		cfg        config
		httpClient *httpclient.Client
	}
)

func newClient(cfg config) *cliClient {
	return &cliClient{
		cfg: cfg,
		httpClient: httpclient.New(httpclient.Config{
			Address:    "http://localhost",
			SocketPath: cfg.SocketPath,
			Timeout:    cfg.Timeout,
		}),
	}

}

func main() {
	cfg, err := bmconfig.GetForCLIClient()
	if err != nil {
		log.Fatal().Err(err).Msg("read config")
	}

	f := parseFlags()
	ctx := context.Background()

	resp, err := newClient(config(*cfg)).process(ctx, f)
	if err != nil {
		log.Fatal().Err(err).Msg("process command")
	}

	log.Info().Any("result", resp).Msg("result")
}

func parseFlags() *flags {
	f := &flags{}

	f.Action = flag.String("action", "", "API action")
	f.ID = flag.Int64("id", 0, "ID")
	f.IDs = flag.String("ids", "", "IDs")
	f.Title = flag.String("title", "", "Title")
	f.Author = flag.String("author", "", "Author")
	f.PublishedDate = flag.String("published_date", "", "Published date (format: '2006-01-02')")
	f.Edition = flag.String("edition", "", "Edition")
	f.Description = flag.String("description", "", "Description")
	f.Genre = flag.String("genre", "", "Genre")
	f.CollectionID = flag.Int64("collection_id", 0, "Collection ID")
	f.BookIDs = flag.String("book_ids", "", "Book IDs")
	f.Name = flag.String("name", "", "Name")
	f.StartDate = flag.String("start_date", "", "Start date (format: '2006-01-02')")
	f.FinishedDate = flag.String("finished_date", "", "Finished date (format: '2006-01-02')")
	f.OrderBy = flag.String("order_by", "", "Order by")
	f.Desc = flag.Bool("desc", false, "Sort in descending order")
	f.Page = flag.Int64("page", 1, "Page number")
	f.PageSize = flag.Int64("page_size", 50, "Page size")

	flag.Parse()

	return f
}

func (c *cliClient) process(ctx context.Context, f *flags) (any, error) {
	switch *f.Action {
	case "get_books":
		return doRequest(ctx, f.toGetBooksReq, c.httpClient.GetBooks)

	case "create_book":
		return doRequest(ctx, f.toCreateBookReq, c.httpClient.CreateBook)

	case "update_book":
		return doRequest(ctx, f.toUpdateBookReq, c.httpClient.UpdateBook)

	case "delete_books":
		return doRequest(ctx, f.toDeleteBooksReq, c.httpClient.DeleteBooks)

	case "get_collections":
		return doRequest(ctx, f.toGetCollectionsReq, c.httpClient.GetCollections)

	case "create_collection":
		return doRequest(ctx, f.toCreateCollectionReq, c.httpClient.CreateCollection)

	case "update_collection":
		return doRequest(ctx, f.toUpdateCollectionReq, c.httpClient.UpdateCollection)

	case "delete_collection":
		return doRequest(ctx, f.toDeleteCollectionReq, c.httpClient.DeleteCollection)

	case "create_books_collection":
		return doRequest(ctx, f.toCreateBooksCollectionReq, c.httpClient.CreateBooksCollection)

	case "delete_books_collection":
		return doRequest(ctx, f.toDeleteBooksCollectionReq, c.httpClient.DeleteBooksCollection)

	default:
		return nil, fmt.Errorf("unknown command: '%s'", *f.Action)
	}
}

func doRequest[Req, Resp any](ctx context.Context, toReq func() (Req, error), doReq func(context.Context, Req) (Resp, error)) (any, error) {
	req, err := toReq()
	if err != nil {
		return nil, fmt.Errorf("convert to request: %w", err)
	}

	resp, err := doReq(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (f *flags) toGetBooksReq() (*api.GetBooksReq, error) {
	var startDate, finishDate time.Time
	var err error

	if f.StartDate != nil {
		startDate, err = parseTime(f.StartDate)
		if err != nil {
			return nil, fmt.Errorf("parse start date: %w", err)
		}
	}

	if f.FinishedDate != nil {
		finishDate, err = parseTime(f.FinishedDate)
		if err != nil {
			return nil, fmt.Errorf("parse finish date: %w", err)
		}
	}

	return &api.GetBooksReq{
		ID:           toValue(f.ID),
		Author:       toValue(f.Author),
		Genre:        toValue(f.Genre),
		CollectionID: toValue(f.CollectionID),
		StartDate:    startDate,
		FinishDate:   finishDate,
		OrderBy:      toValue(f.OrderBy),
		Desc:         toValue(f.Desc),
		Page:         toValue(f.Page),
		PageSize:     toValue(f.PageSize),
	}, nil
}

func (f *flags) toCreateBookReq() (*api.CreateBookReq, error) {
	publishedDate, err := parseTime(f.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("parse published date: %w", err)
	}

	return &api.CreateBookReq{
		Title:         toValue(f.Title),
		Author:        toValue(f.Author),
		PublishedDate: publishedDate,
		Edition:       toValue(f.Edition),
		Description:   toValue(f.Description),
		Genre:         toValue(f.Genre),
	}, nil
}

func (f *flags) toUpdateBookReq() (*api.UpdateBookReq, error) {
	var publishedDate time.Time
	var err error
	if f.PublishedDate != nil && *f.PublishedDate != "" {
		publishedDate, err = parseTime(f.PublishedDate)
		if err != nil {
			return nil, fmt.Errorf("parse published date: %w", err)
		}
	}

	return &api.UpdateBookReq{
		ID:            toValue(f.ID),
		Title:         toValue(f.Title),
		Author:        toValue(f.Author),
		PublishedDate: publishedDate,
		Edition:       toValue(f.Edition),
		Description:   toValue(f.Description),
		Genre:         toValue(f.Genre),
	}, nil
}

func (f *flags) toDeleteBooksReq() (*api.DeleteBooksReq, error) {
	ids := make([]int64, 0)
	for _, idStr := range strings.Split(*f.IDs, ",") {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("incorrect book id: '%s'", idStr)
		}

		ids = append(ids, id)
	}

	return &api.DeleteBooksReq{
		IDs: ids,
	}, nil
}

func (f *flags) toGetCollectionsReq() (*api.GetCollectionsReq, error) {
	ids := make([]int64, 0)
	for _, idStr := range strings.Split(*f.IDs, ",") {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("incorrect book id: '%s'", idStr)
		}

		ids = append(ids, id)
	}

	return &api.GetCollectionsReq{
		IDs:      ids,
		OrderBy:  toValue(f.OrderBy),
		Desc:     toValue(f.Desc),
		Page:     toValue(f.Page),
		PageSize: toValue(f.PageSize),
	}, nil
}

func (f *flags) toCreateCollectionReq() (*api.CreateCollectionReq, error) {
	return &api.CreateCollectionReq{
		Name:        toValue(f.Name),
		Description: toValue(f.Description),
	}, nil
}

func (f *flags) toUpdateCollectionReq() (*api.UpdateCollectionReq, error) {
	return &api.UpdateCollectionReq{
		ID:          toValue(f.ID),
		Name:        toValue(f.Name),
		Description: toValue(f.Description),
	}, nil
}

func (f *flags) toDeleteCollectionReq() (*api.DeleteCollectionReq, error) {
	id := toValue(f.ID)
	if id == 0 {
		return nil, fmt.Errorf("collection id is required")
	}

	return &api.DeleteCollectionReq{
		ID: id,
	}, nil
}

func (f *flags) toCreateBooksCollectionReq() (*api.CreateBooksCollectionReq, error) {
	bookIDs := make([]int64, 0)
	for _, idStr := range strings.Split(*f.BookIDs, ",") {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("incorrect book id: '%s'", idStr)
		}

		bookIDs = append(bookIDs, id)
	}

	cid := toValue(f.CollectionID)
	if cid <= 0 || len(bookIDs) == 0 {
		return nil, fmt.Errorf("CollectionID and BookID are required for create_books_collection action")
	}

	return &api.CreateBooksCollectionReq{
		CID:     cid,
		BookIDs: bookIDs,
	}, nil
}

func (f *flags) toDeleteBooksCollectionReq() (*api.DeleteBooksCollectionReq, error) {
	bookIDs := make([]int64, 0)
	for _, idStr := range strings.Split(*f.BookIDs, ",") {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("incorrect book id: '%s'", idStr)
		}

		bookIDs = append(bookIDs, id)
	}

	if f.CollectionID == nil || len(bookIDs) == 0 {
		return nil, fmt.Errorf("both of collection_id and book_ids are required")
	}

	return &api.DeleteBooksCollectionReq{
		CID:     toValue(f.CollectionID),
		BookIDs: bookIDs,
	}, nil
}

func parseTime(dateStr *string) (time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return time.Time{}, nil
	}

	parsedTime, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func toValue[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}

	var v T

	return v
}
