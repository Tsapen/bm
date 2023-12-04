package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	bmconfig "github.com/Tsapen/bm/internal/config"
	"github.com/Tsapen/bm/pkg/api"
	httpclient "github.com/Tsapen/bm/pkg/http-client"
)

const (
	formatDate = "2006-01-02"
)

type (
	flags struct {
		Action        *string
		ID            *int64
		IDs           *ListValue
		Title         *string
		Author        *string
		PublishedDate *DateValue
		Edition       *string
		Description   *string
		Genre         *string
		CollectionID  *int64
		BookIDs       *ListValue
		Name          *string
		StartDate     *DateValue
		FinishedDate  *DateValue
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
		fmt.Fprintf(os.Stderr, "read config: %v\n", err)
		os.Exit(1)
	}

	f := parseFlags()
	ctx := context.Background()

	resp, err := newClient(config(*cfg)).process(ctx, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "process command: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, resp)
	os.Exit(0)
}

func parseFlags() *flags {
	f := &flags{
		IDs:          new(ListValue),
		BookIDs:      new(ListValue),
		StartDate:    new(DateValue),
		FinishedDate: new(DateValue),
	}

	f.Action = flag.String("action", "", "API action")
	f.ID = flag.Int64("id", 0, "ID")
	flag.Var(f.IDs, "ids", "IDs")
	f.Title = flag.String("title", "", "Title")
	f.Author = flag.String("author", "", "Author")
	flag.Var(f.PublishedDate, "date", "Published date (format: '2006-01-02')")
	f.Edition = flag.String("edition", "", "Edition")
	f.Description = flag.String("description", "", "Description")
	f.Genre = flag.String("genre", "", "Genre")
	f.CollectionID = flag.Int64("collection_id", 0, "Collection ID")
	flag.Var(f.BookIDs, "book_ids", "Book IDs")
	f.Name = flag.String("name", "", "Name")
	flag.Var(f.StartDate, "start_date", "Start date (format: '2006-01-02')")
	flag.Var(f.FinishedDate, "finished_date", "Finished date (format: '2006-01-02')")
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
	return &api.GetBooksReq{
		Author:       toValue(f.Author),
		Genre:        toValue(f.Genre),
		CollectionID: toValue(f.CollectionID),
		StartDate:    f.StartDate.date(),
		FinishDate:   f.FinishedDate.Date,
		OrderBy:      toValue(f.OrderBy),
		Desc:         toValue(f.Desc),
		Page:         toValue(f.Page),
		PageSize:     toValue(f.PageSize),
	}, nil
}

func (f *flags) toCreateBookReq() (*api.CreateBookReq, error) {
	return &api.CreateBookReq{
		Title:         toValue(f.Title),
		Author:        toValue(f.Author),
		PublishedDate: f.PublishedDate.date(),
		Edition:       toValue(f.Edition),
		Description:   toValue(f.Description),
		Genre:         toValue(f.Genre),
	}, nil
}

func (f *flags) toUpdateBookReq() (*api.UpdateBookReq, error) {
	return &api.UpdateBookReq{
		ID:            toValue(f.ID),
		Title:         toValue(f.Title),
		Author:        toValue(f.Author),
		PublishedDate: f.PublishedDate.date(),
		Edition:       toValue(f.Edition),
		Description:   toValue(f.Description),
		Genre:         toValue(f.Genre),
	}, nil
}

func (f *flags) toDeleteBooksReq() (*api.DeleteBooksReq, error) {
	return &api.DeleteBooksReq{
		IDs: f.IDs.IDs,
	}, nil
}

func (f *flags) toGetCollectionsReq() (*api.GetCollectionsReq, error) {
	return &api.GetCollectionsReq{
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
	cid := toValue(f.CollectionID)
	if cid <= 0 || len(f.BookIDs.IDs) == 0 {
		return nil, fmt.Errorf("CollectionID and BookID are required for create_books_collection action")
	}

	return &api.CreateBooksCollectionReq{
		CID:     cid,
		BookIDs: f.BookIDs.IDs,
	}, nil
}

func (f *flags) toDeleteBooksCollectionReq() (*api.DeleteBooksCollectionReq, error) {
	if f.CollectionID == nil || len(f.BookIDs.IDs) == 0 {
		return nil, fmt.Errorf("both of collection_id and book_ids are required")
	}

	return &api.DeleteBooksCollectionReq{
		CID:     toValue(f.CollectionID),
		BookIDs: f.BookIDs.IDs,
	}, nil
}

type DateValue struct {
	Date time.Time
}

func (v *DateValue) date() time.Time {
	if v != nil {
		return v.Date
	}

	return time.Time{}
}

func (v *DateValue) String() string {
	if v != nil {
		return v.Date.Format(formatDate)
	}
	return ""
}

func (v *DateValue) Set(str string) error {
	t, err := time.Parse(formatDate, str)
	if err != nil {
		return err
	}

	v.Date = t

	return nil
}

type ListValue struct {
	IDs []int64
}

func (v *ListValue) String() string {
	if v.IDs != nil {
		return fmt.Sprint(v.IDs)
	}
	return ""
}

func (v *ListValue) Set(str string) error {
	if str == "" {
		return nil
	}

	ids := make([]int64, 0)
	for _, idStr := range strings.Split(str, ",") {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return fmt.Errorf("incorrect id: '%s'", idStr)
		}

		ids = append(ids, id)
	}

	v.IDs = ids

	return nil
}

func toValue[T any](ptr *T) (v T) {
	if ptr != nil {
		return *ptr
	}

	return
}
