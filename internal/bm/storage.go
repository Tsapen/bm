package bm

import (
	"context"
	"time"
)

type (
	BookFilter struct {
		ID           int64
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
	// Books retrieves a list of books based on the provided filter criteria.
	Books(ctx context.Context, f BookFilter) ([]Book, error)

	// CreateBook creates a new book with the provided details.
	CreateBook(ctx context.Context, b Book) (id int64, err error)

	// UpdateBook updates an existing book with the provided details.
	UpdateBook(ctx context.Context, b Book) error

	// DeleteBooks deletes books based on their IDs.
	DeleteBooks(ctx context.Context, ids []int64) error

	// Collections retrieves a list of collections based on the provided filter criteria.
	Collections(ctx context.Context, f CollectionsFilter) ([]Collection, error)

	// CreateCollection creates a new collection with the provided details.
	CreateCollection(ctx context.Context, c Collection) (int64, error)

	// UpdateCollection updates an existing collection with the provided details.
	UpdateCollection(ctx context.Context, c Collection) error

	// DeleteCollection deletes a collection based on its ID.
	DeleteCollection(ctx context.Context, id int64) error

	// CreateBooksCollection adds a list of books to an existing collection.
	CreateBooksCollection(ctx context.Context, collectionID int64, bookIDs []int64) error

	// DeleteBooksCollection removes a list of books from an existing collection.
	DeleteBooksCollection(ctx context.Context, collectionID int64, bookIDs []int64) error
}
