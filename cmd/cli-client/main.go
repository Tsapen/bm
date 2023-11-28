package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	cliclient "github.com/Tsapen/bm/internal/cli-client"
	"github.com/Tsapen/bm/internal/config"
	"github.com/Tsapen/bm/pkg/api"
)

type (
	flags struct {
		Action        *string
		ID            *int64
		Title         *string
		Author        *string
		PublishedDate *string
		Edition       *string
		Description   *string
		Genre         *string
		CollectionID  *int64
		CollectionsID *string
		BooksID       *string
		Name          *string
		StartDate     *string
		FinishedDate  *string
		OrderBy       *string
		Desc          *bool
		Page          *int64
		PageSize      *int64
	}

	command struct {
		Action string `json:"action"`
		Data   any    `json:"data"`
	}
)

func main() {
	cfg, err := config.GetForCLIClient()
	if err != nil {
		log.Fatal().Err(err).Msg("read config")
	}

	cmd, err := getCommand()
	if err != nil {
		log.Fatal().Err(err).Msg("get command")
	}

	log.Info().Any("command", cmd).Msg("cmd")
	resp, err := cliclient.New(cliclient.Config(*cfg.UnixSocketCfg)).DoRequest(cmd)
	if err != nil {
		log.Fatal().Err(err).Msg("do request to unix-socket server")
	}

	log.Info().Any("result", resp)
}

func parseFlags() *flags {
	f := &flags{}

	f.Action = flag.String("action", "", "API action")
	f.ID = flag.Int64("id", 0, "ID")
	f.Title = flag.String("title", "", "Title")
	f.Author = flag.String("author", "", "Author")
	f.PublishedDate = flag.String("published_date", "", "Published date (format: '2006-01-02')")
	f.Edition = flag.String("edition", "", "Edition")
	f.Description = flag.String("description", "", "Description")
	f.Genre = flag.String("genre", "", "Genre")
	f.CollectionID = flag.Int64("collection_id", 0, "Collection ID")
	f.CollectionsID = flag.String("collections_id", "", "Collections ID")
	f.BooksID = flag.String("books_id", "", "Book IDs")
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

func getCommand() (*command, error) {
	f := parseFlags()

	action := *f.Action
	var data any
	var err error
	switch action {
	case "get_books":
		data, err = f.toGetBooksReq()

	case "create_book":
		data, err = f.toCreateBookReq()

	case "update_book":
		data, err = f.toUpdateBookReq()

	case "delete_book":
		data, err = f.toDeleteBooksReq()

	case "get_collections":
		data, err = f.toGetCollectionsReq()

	case "create_collection":
		data, err = f.toCreateCollectionReq()

	case "update_collection":
		data, err = f.toUpdateCollectionReq()

	case "delete_collection":
		data, err = f.toDeleteCollectionReq()

	case "create_books_collection":
		data, err = f.toCreateBooksCollectionReq()

	case "delete_books_collection":
		data, err = f.toDeleteBooksCollectionReq()

	default:
		return nil, fmt.Errorf("unknown command: '%s'", action)
	}

	if err != nil {
		return nil, fmt.Errorf("construct command: %w", err)
	}

	return &command{
		Action: action,
		Data:   data,
	}, nil
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
	publishedDate, err := parseTime(f.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("parse published date: %w", err)
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
	if f.ID == nil {
		return nil, fmt.Errorf("ID is required for delete_books action")
	}

	return &api.DeleteBooksReq{
		IDs: []int64{*f.ID},
	}, nil
}

func (f *flags) toGetCollectionsReq() (*api.GetCollectionsReq, error) {
	ids, err := toArray(toValue(f.CollectionsID))
	if err != nil {
		return nil, fmt.Errorf("parse ids: %w", err)
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
	cid := toValue(f.CollectionID)
	if cid == 0 {
		return nil, fmt.Errorf("collection id is required")
	}

	return &api.DeleteCollectionReq{
		ID: cid,
	}, nil
}

func (f *flags) toCreateBooksCollectionReq() (*api.CreateBooksCollectionReq, error) {
	if f.CollectionID == nil || f.ID == nil {
		return nil, fmt.Errorf("CollectionID and BookID are required for create_books_collection action")
	}

	return &api.CreateBooksCollectionReq{
		CID:     toValue(f.CollectionID),
		BookIDs: []int64{*f.ID},
	}, nil
}

func (f *flags) toDeleteBooksCollectionReq() (*api.DeleteBooksCollectionReq, error) {
	bookIDs, err := toArray(toValue(f.BooksID))
	if err != nil {
		return nil, fmt.Errorf("book ids are required")
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

func toArray(ids string) ([]int64, error) {
	values := strings.Split(ids, ",")

	result := make([]int64, 0, len(values))
	for _, v := range values {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("convert '%s' to int64: %w", v, err)
		}

		result = append(result, id)
	}

	return result, nil
}

func toValue[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}

	var v T

	return v
}
