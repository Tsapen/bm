package bm

import (
	"context"
	"time"
)

type (
	BookFilter struct {
		Author       string
		Genre        string
		CollectionID int64
		StartDate    time.Time
		FinishDate   time.Time
		OrderBy      string
		Desc         bool
		Page         int64
		PageSize     int64
	}

	Book struct {
		ID            int64     `db:"id"`
		Title         string    `db:"title"`
		Author        string    `db:"author"`
		PublishedDate time.Time `db:"published_date"`
		Edition       string    `db:"edition"`
		Description   string    `db:"description"`
		Genre         string    `db:"genre"`
	}

	Collection struct {
		ID          int64  `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
	}

	CollectionInfo struct {
		Collection
		Books []Book
	}

	CollectionsFilter struct {
		IDs      []int64
		OrderBy  string
		Desc     bool
		Page     int64
		PageSize int64
	}

	BooksCollectionFilter struct {
		CID      int64
		OrderBy  string
		Desc     bool
		Page     int64
		PageSize int64
	}
)

// Storage is a database interface.
type Storage interface {
	Books(context.Context, BookFilter) ([]Book, error)
	CreateBook(context.Context, Book) (id int64, err error)
	UpdateBook(context.Context, Book) error
	DeleteBooks(ctx context.Context, ids []int64) error

	Collections(context.Context, CollectionsFilter) ([]Collection, error)
	CreateCollection(context.Context, Collection) (int64, error)
	UpdateCollection(context.Context, Collection) error
	DeleteCollection(context.Context, int64) error

	CreateBooksCollection(ctx context.Context, collectionID int64, bookIDs []int64) error
	DeleteBooksCollection(ctx context.Context, collectionID int64, bookIDs []int64) error
}
