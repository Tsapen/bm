package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	bm "github.com/Tsapen/bm/internal/bm"
)

// Collection gets collection by its id.
func (s *DB) Collection(ctx context.Context, id int64) (*bm.Collection, error) {
	q := `SELECT id, name, description FROM collections c 
			WHERE id=$1
	`

	collection := new(bm.Collection)
	err := s.GetContext(ctx, collection, q, id)
	if err != nil {
		return nil, bm.NewInternalError("select book: %w", err)
	}

	return collection, nil
}

// Collections gets collections.
func (s *DB) Collections(ctx context.Context, f bm.CollectionsFilter) ([]bm.Collection, error) {
	q := "SELECT c.id, c.name, c.description FROM collections c "
	q += pagination(f.Page, f.PageSize)

	var collections []bm.Collection
	if err := s.SelectContext(ctx, &collections, q); err != nil {
		return nil, bm.NewInternalError("select collections: %w", err)
	}

	return collections, nil
}

func (s *DB) CreateCollection(ctx context.Context, c bm.Collection) (int64, error) {
	query := `
		INSERT INTO collections (name, description)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	err := s.QueryRowContext(ctx, query, c.Name, c.Description).Scan(&id)
	if err != nil {
		return 0, bm.NewInternalError("insert collection: %w", err)
	}

	return id, nil
}

// UpdateCollection updates a collection and its books.
func (s *DB) UpdateCollection(ctx context.Context, c bm.Collection) (err error) {
	params := []any{c.Name, c.Description, c.ID}
	q := `UPDATE collections c SET
			name = $1,
			description = $2
		WHERE id = $3`

	result, err := s.DB.ExecContext(ctx, q, params...)
	if err != nil {
		return bm.NewInternalError("update collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return bm.NewInternalError("get the number of affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return bm.NewNotFoundError("collection with ID %d not found", c.ID)
	}

	return nil
}

// DeleteCollections deletes collection and its associations.
func (s *DB) DeleteCollection(ctx context.Context, id int64) error {
	err := s.withTX(ctx, func(tx *sql.Tx) error {
		q := `DELETE FROM books_collection bc WHERE bc.collection_id = $1`
		if _, err := tx.ExecContext(ctx, q, id); err != nil {
			return bm.NewInternalError("delete collection books: %w", err)
		}

		q = `DELETE FROM collections WHERE id = $1`
		if _, err := tx.ExecContext(ctx, q, id); err != nil {
			return bm.NewInternalError("delete collection: %w", err)
		}

		return nil
	})
	if err != nil {
		return bm.NewInternalError("execute tx: %w", err)
	}

	return nil
}

func insertBooksCollectionValues(numValues int) string {
	values := make([]string, 0, numValues)
	for i := 0; i < numValues; i++ {
		values = append(values, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
	}

	return strings.Join(values, ", ")
}

// CreateBooksCollection adds books to a collection.
func (s *DB) CreateBooksCollection(ctx context.Context, cID int64, bookIDs []int64) error {
	err := s.withTX(ctx, func(tx *sql.Tx) error {
		q := fmt.Sprintf(
			`INSERT INTO books_collection (collection_id, book_id) VALUES %s`,
			insertBooksCollectionValues(len(bookIDs)))

		args := make([]interface{}, 0, len(bookIDs)*2)
		for _, bookID := range bookIDs {
			args = append(args, cID, bookID)
		}

		if _, err := tx.ExecContext(ctx, q, args...); err != nil {
			return bm.NewInternalError("add books to collection: %w", err)
		}

		return nil
	})
	if err != nil {
		return bm.NewInternalError("execute tx: %w", err)
	}

	return nil
}

// DeleteBooksCollection deletes books from a collection.
func (s *DB) DeleteBooksCollection(ctx context.Context, cID int64, bookIDs []int64) error {
	q := `DELETE FROM books_collection bc WHERE bc.collection_id = $1 AND bc.book_id = ANY ($2)`
	result, err := s.ExecContext(ctx, q, cID, pq.Array(bookIDs))
	if err != nil {
		return bm.NewInternalError("remove books from collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return bm.NewInternalError("get the number of affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return bm.NewNotFoundError("book with ID %v not found in collection %d", bookIDs, cID)
	}

	return nil
}
