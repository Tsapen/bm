package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	bmconfig "github.com/Tsapen/bm/internal/config"
	"github.com/Tsapen/bm/pkg/api"
	httpclient "github.com/Tsapen/bm/pkg/http-client"
)

const (
	formatDate = "2006-01-02"
)

func process[Req, Resp any](
	ctx context.Context,
	toReq func() (Req, error),
	doReq func(context.Context, Req) (Resp, error),
) {
	req, err := toReq()
	if err != nil {
		fmt.Fprintf(os.Stderr, "convert to api request: %v\n", err)
		os.Exit(1)
	}

	resp, err := doReq(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "do request: %v\n", err)
		os.Exit(1)
	}

	var output any = resp
	if _, ok := output.(bool); ok {
		output = map[string]any{
			"success": output,
		}
	}
	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "print response: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func main() {
	cfg, err := bmconfig.GetForCLIClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "read config: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	c := newClient(config(*cfg))

	getBookReq := new(getBookReqCli)
	var cmdGetBook = &cobra.Command{
		Use:   "get_book [id to book]",
		Short: "returns json with book data",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, getBookReq.toAPIReq, c.httpClient.GetBook)
		},
	}

	cmdGetBook.Flags().Int64Var(&getBookReq.ID, "id", 0, "Source directory to read from")
	cmdGetBook.MarkFlagRequired("id")

	createBookReq := new(createBookReqCli)
	cmdCreateBook := &cobra.Command{
		Use:   "create_book",
		Short: "Creates a new book",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, createBookReq.toAPIReq, c.httpClient.CreateBook)
		},
	}

	cmdCreateBook.Flags().StringVarP(&createBookReq.Title, "title", "t", "", "Title of the book (required)")
	cmdCreateBook.Flags().StringVarP(&createBookReq.Author, "author", "a", "", "Author of the book (required)")
	cmdCreateBook.Flags().StringVarP(&createBookReq.PublishedDate, "published_date", "d", "", "Published date of the book in the format YYYY-MM-DD (required)")
	cmdCreateBook.Flags().StringVar(&createBookReq.Edition, "edition", "", "Edition of the book")
	cmdCreateBook.Flags().StringVar(&createBookReq.Description, "description", "", "Description of the book")
	cmdCreateBook.Flags().StringVar(&createBookReq.Genre, "genre", "", "Genre of the book")

	cmdCreateBook.MarkFlagRequired("title")
	cmdCreateBook.MarkFlagRequired("author")
	cmdCreateBook.MarkFlagRequired("published_date")

	getBooksReq := new(getBooksReqCli)
	cmdGetBooks := &cobra.Command{
		Use:   "get_books",
		Short: "Get a list of books",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, getBooksReq.toAPIReq, c.httpClient.GetBooks)
		},
	}

	cmdGetBooks.Flags().StringVar(&getBooksReq.Author, "author", "", "Author of the books")
	cmdGetBooks.Flags().StringVar(&getBooksReq.Genre, "genre", "", "Genre of the books")
	cmdGetBooks.Flags().Int64Var(&getBooksReq.CollectionID, "collection_id", 0, "ID of the collection")
	cmdGetBooks.Flags().StringVar(&getBooksReq.StartDate, "start_date", "", "Start date in the format YYYY-MM-DD")
	cmdGetBooks.Flags().StringVar(&getBooksReq.FinishDate, "finish_date", "", "Finish date in the format YYYY-MM-DD")
	cmdGetBooks.Flags().StringVar(&getBooksReq.OrderBy, "order_by", "", "Order by a specific field")
	cmdGetBooks.Flags().BoolVar(&getBooksReq.Desc, "desc", false, "Sort in descending order")
	cmdGetBooks.Flags().Int64Var(&getBooksReq.Page, "page", 1, "Page number")
	cmdGetBooks.Flags().Int64Var(&getBooksReq.PageSize, "page_size", 10, "Number of items per page")

	updateBooksReq := new(updateBookReqCli)
	var cmdUpdateBooks = &cobra.Command{
		Use:   "update_book",
		Short: "Update existing books",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, updateBooksReq.toAPIReq, c.httpClient.UpdateBook)
		},
	}

	cmdUpdateBooks.Flags().Int64Var(&updateBooksReq.ID, "id", 0, "ID of the book to update (required)")
	cmdUpdateBooks.Flags().StringVar(&updateBooksReq.Title, "title", "", "Updated title of the book")
	cmdUpdateBooks.Flags().StringVar(&updateBooksReq.Author, "author", "", "Updated author of the book")
	cmdUpdateBooks.Flags().StringVar(&updateBooksReq.PublishedDate, "published_date", "", "Updated published date of the book in the format YYYY-MM-DD")
	cmdUpdateBooks.Flags().StringVar(&updateBooksReq.Edition, "edition", "", "Updated edition of the book")
	cmdUpdateBooks.Flags().StringVar(&updateBooksReq.Description, "description", "", "Updated description of the book")
	cmdUpdateBooks.Flags().StringVar(&updateBooksReq.Genre, "genre", "", "Updated genre of the book")
	cmdUpdateBooks.MarkFlagRequired("id")

	var deleteBooksReq = &deleteBooksReqCli{}
	var cmdDeleteBooks = &cobra.Command{
		Use:   "delete_books",
		Short: "Delete existing books",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, deleteBooksReq.toAPIReq, c.httpClient.DeleteBooks)
		},
	}

	cmdDeleteBooks.Flags().Int64SliceVar(&deleteBooksReq.IDs, "ids", nil, "IDs of the books to delete (comma-separated) (required)")
	cmdDeleteBooks.MarkFlagRequired("ids")

	var getCollectionReq = &getCollectionReqCli{}
	var getCollectionsReq = &getCollectionsReqCli{}

	var cmdGetCollection = &cobra.Command{
		Use:   "get_collection",
		Short: "Get information about a collection",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, getCollectionReq.toAPIReq, c.httpClient.GetCollection)
		},
	}

	cmdGetCollection.Flags().Int64Var(&getCollectionReq.ID, "id", 0, "ID of the collection to retrieve (required)")
	cmdGetCollection.MarkFlagRequired("id")

	var cmdGetCollections = &cobra.Command{
		Use:   "get_collections",
		Short: "Get a list of collections",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, getCollectionsReq.toAPIReq, c.httpClient.GetCollections)
		},
	}

	cmdGetCollections.Flags().StringVar(&getCollectionsReq.OrderBy, "order_by", "", "Order by a specific field")
	cmdGetCollections.Flags().BoolVar(&getCollectionsReq.Desc, "desc", false, "Sort in descending order")
	cmdGetCollections.Flags().Int64Var(&getCollectionsReq.Page, "page", 1, "Page number")
	cmdGetCollections.Flags().Int64Var(&getCollectionsReq.PageSize, "page_size", 10, "Number of items per page")

	var createCollectionReq = &createCollectionReqCli{}

	var cmdCreateCollection = &cobra.Command{
		Use:   "create_collection",
		Short: "Create a new collection",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, createCollectionReq.toAPIReq, c.httpClient.CreateCollection)
		},
	}

	cmdCreateCollection.Flags().StringVar(&createCollectionReq.Name, "name", "", "Name of the new collection (required)")
	cmdCreateCollection.Flags().StringVar(&createCollectionReq.Description, "description", "", "Description of the new collection")

	cmdCreateCollection.MarkFlagRequired("name")

	var updateCollectionReq = &updateCollectionReqCli{}
	var cmdUpdateCollection = &cobra.Command{
		Use:   "update_collection",
		Short: "Update an existing collection",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, updateCollectionReq.toAPIReq, c.httpClient.UpdateCollection)
		},
	}

	cmdUpdateCollection.Flags().Int64Var(&updateCollectionReq.ID, "id", 0, "ID of the collection to update (required)")
	cmdUpdateCollection.Flags().StringVar(&updateCollectionReq.Name, "name", "", "Updated name of the collection")
	cmdUpdateCollection.Flags().StringVar(&updateCollectionReq.Description, "description", "", "Updated description of the collection")

	cmdUpdateCollection.MarkFlagRequired("id")

	var deleteCollectionsReq = &deleteCollectionsReqCli{}
	var cmdDeleteCollections = &cobra.Command{
		Use:   "delete_collection",
		Short: "Delete existing collections",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, deleteCollectionsReq.toAPIReq, c.httpClient.DeleteCollection)
		},
	}

	cmdDeleteCollections.Flags().Int64Var(&deleteCollectionsReq.ID, "id", 0, "ID of the collection to delete (required)")
	cmdDeleteCollections.MarkFlagRequired("ids")

	var createBooksCollectionReq = &createBooksCollectionReqCli{}

	var cmdCreateBooksCollection = &cobra.Command{
		Use:   "create_books_collection",
		Short: "Create a new association between books and a collection",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, createBooksCollectionReq.toAPIReq, c.httpClient.CreateBooksCollection)
		},
	}

	cmdCreateBooksCollection.Flags().Int64Var(&createBooksCollectionReq.CID, "collection_id", 0, "ID of the collection to associate books with (required)")
	cmdCreateBooksCollection.Flags().Int64SliceVar(&createBooksCollectionReq.BookIDs, "book_ids", nil, "IDs of the books to associate with the collection (comma-separated) (required)")
	cmdCreateBooksCollection.MarkFlagRequired("collection_id")
	cmdCreateBooksCollection.MarkFlagRequired("book_ids")

	var deleteBooksCollectionReq = &deleteBooksCollectionReqCli{}
	var cmdDeleteBooksCollection = &cobra.Command{
		Use:   "delete_books_collection",
		Short: "Delete an association between books and a collection",
		Run: func(cmd *cobra.Command, args []string) {
			process(ctx, deleteBooksCollectionReq.toAPIReq, c.httpClient.DeleteBooksCollection)
		},
	}

	cmdDeleteBooksCollection.Flags().Int64Var(&deleteBooksCollectionReq.CID, "collection_id", 0, "ID of the collection to disassociate books from (required)")
	cmdDeleteBooksCollection.Flags().Int64SliceVar(&deleteBooksCollectionReq.BookIDs, "book_ids", nil, "IDs of the books to disassociate from the collection (comma-separated) (required)")
	cmdDeleteBooksCollection.MarkFlagRequired("collection_id")
	cmdDeleteBooksCollection.MarkFlagRequired("book_ids")

	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(
		cmdGetBook,
		cmdCreateBook,
		cmdGetBooks,
		cmdUpdateBooks,
		cmdDeleteBooks,
		cmdGetCollection,
		cmdGetCollections,
		cmdCreateCollection,
		cmdUpdateCollection,
		cmdDeleteCollections,
		cmdCreateBooksCollection,
		cmdDeleteBooksCollection,
	)
	rootCmd.Execute()
}

type (
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

type getBookReqCli struct {
	ID int64
}

func (r *getBookReqCli) toAPIReq() (*api.GetBookReq, error) {
	return &api.GetBookReq{
		ID: r.ID,
	}, nil
}

type getBooksReqCli struct {
	Author       string
	Genre        string
	CollectionID int64
	StartDate    string
	FinishDate   string
	OrderBy      string
	Desc         bool
	Page         int64
	PageSize     int64
}

func (r *getBooksReqCli) toAPIReq() (*api.GetBooksReq, error) {
	req := &api.GetBooksReq{
		Author:       r.Author,
		Genre:        r.Genre,
		CollectionID: r.CollectionID,
		OrderBy:      r.OrderBy,
		Desc:         r.Desc,
		Page:         r.Page,
		PageSize:     r.PageSize,
	}

	var err error

	if r.StartDate != "" {
		req.StartDate, err = time.Parse(formatDate, r.StartDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse start_date: %v", err)
		}
	}

	if r.FinishDate != "" {
		req.FinishDate, err = time.Parse(formatDate, r.FinishDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse finish_date: %v", err)
		}
	}

	return req, nil
}

type createBookReqCli struct {
	Title         string
	Author        string
	PublishedDate string
	Edition       string
	Description   string
	Genre         string
}

func (r *createBookReqCli) toAPIReq() (*api.CreateBookReq, error) {
	parsedDate, err := time.Parse(formatDate, r.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse published_date: %v", err)
	}

	return &api.CreateBookReq{
		Title:         r.Title,
		Author:        r.Author,
		PublishedDate: parsedDate,
		Edition:       r.Edition,
		Description:   r.Description,
		Genre:         r.Genre,
	}, nil
}

type updateBookReqCli struct {
	ID            int64
	Title         string
	Author        string
	PublishedDate string
	Edition       string
	Description   string
	Genre         string
}

type deleteBooksReqCli struct {
	IDs []int64
}

func (r *updateBookReqCli) toAPIReq() (*api.UpdateBookReq, error) {
	var publishedDate time.Time
	var err error

	if r.PublishedDate != "" {
		publishedDate, err = time.Parse(formatDate, r.PublishedDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse published_date: %v", err)
		}
	}

	return &api.UpdateBookReq{
		ID:            r.ID,
		Title:         r.Title,
		Author:        r.Author,
		PublishedDate: publishedDate,
		Edition:       r.Edition,
		Description:   r.Description,
		Genre:         r.Genre,
	}, nil
}

func (r *deleteBooksReqCli) toAPIReq() (*api.DeleteBooksReq, error) {
	return &api.DeleteBooksReq{
		IDs: r.IDs,
	}, nil
}

type getCollectionReqCli struct {
	ID int64
}

type getCollectionsReqCli struct {
	OrderBy  string
	Desc     bool
	Page     int64
	PageSize int64
}

func (r *getCollectionReqCli) toAPIReq() (*api.GetCollectionReq, error) {
	return &api.GetCollectionReq{
		ID: r.ID,
	}, nil
}

func (r *getCollectionsReqCli) toAPIReq() (*api.GetCollectionsReq, error) {
	return &api.GetCollectionsReq{
		OrderBy:  r.OrderBy,
		Desc:     r.Desc,
		Page:     r.Page,
		PageSize: r.PageSize,
	}, nil
}

type createCollectionReqCli struct {
	Name        string
	Description string
}

func (r *createCollectionReqCli) toAPIReq() (*api.CreateCollectionReq, error) {
	return &api.CreateCollectionReq{
		Name:        r.Name,
		Description: r.Description,
	}, nil
}

type updateCollectionReqCli struct {
	ID          int64
	Name        string
	Description string
}

func (r *updateCollectionReqCli) toAPIReq() (*api.UpdateCollectionReq, error) {
	return &api.UpdateCollectionReq{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}, nil
}

type deleteCollectionsReqCli struct {
	ID int64
}

func (r *deleteCollectionsReqCli) toAPIReq() (*api.DeleteCollectionReq, error) {
	return &api.DeleteCollectionReq{
		ID: r.ID,
	}, nil
}

type createBooksCollectionReqCli struct {
	CID     int64
	BookIDs []int64
}

func (r *createBooksCollectionReqCli) toAPIReq() (*api.CreateBooksCollectionReq, error) {
	return &api.CreateBooksCollectionReq{
		CID:     r.CID,
		BookIDs: r.BookIDs,
	}, nil
}

type deleteBooksCollectionReqCli struct {
	CID     int64
	BookIDs []int64
}

func (r *deleteBooksCollectionReqCli) toAPIReq() (*api.DeleteBooksCollectionReq, error) {
	return &api.DeleteBooksCollectionReq{
		CID:     r.CID,
		BookIDs: r.BookIDs,
	}, nil
}
