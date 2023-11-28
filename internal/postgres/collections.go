package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	bm "github.com/Tsapen/bm/internal/bm"
)

// Collections gets collections.
func (s *DB) Collections(ctx context.Context, f bm.CollectionsFilter) ([]bm.Collection, error) {
	q := "SELECT c.id, c.name, c.description FROM collections c "
	whereClause, params := collectionsWhereClause(f)
	q += whereClause
	q += orderBy("c", f.OrderBy, f.Desc)
	q += pagination(f.Page, f.PageSize)

	rows, err := s.NamedQueryContext(ctx, q, params)
	if err != nil {
		return nil, fmt.Errorf("select collections: %w", err)
	}

	defer func() {
		err = bm.HandleErrPair(rows.Close(), err)
	}()

	var collections []bm.Collection
	if err = sqlx.StructScan(rows, &collections); err != nil {
		return nil, fmt.Errorf("copy data into struct: %w", err)
	}

	return collections, nil
}

func collectionsWhereClause(f bm.CollectionsFilter) (string, map[string]any) {
	whereClauses := make([]string, 0)
	params := make(map[string]any, 0)

	if len(f.IDs) > 0 {
		whereClauses = append(whereClauses, "c.id = ANY(:ids) ")
		params["ids"] = pq.Array(f.IDs)
	}

	if len(whereClauses) == 0 {
		return "", map[string]any{}
	}

	return "WHERE " + strings.Join(whereClauses, " AND "), params
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
		return 0, fmt.Errorf("execute collection insert: %w", err)
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
		return fmt.Errorf("update collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get the number of affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("collection with ID %d not found", c.ID)
	}

	return nil
}

// DeleteCollections deletes collection and its associations.
func (s *DB) DeleteCollection(ctx context.Context, id int64) error {
	err := s.withTX(ctx, func(tx *sql.Tx) error {
		q := `DELETE FROM books_collection bc WHERE bc.collection_id = $1`
		if _, err := tx.ExecContext(ctx, q, id); err != nil {
			return fmt.Errorf("delete collection books: %w", err)
		}

		q = `DELETE FROM collections WHERE id = $1`
		if _, err := tx.ExecContext(ctx, q, id); err != nil {
			return fmt.Errorf("delete collection: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("execute tx: %w", err)
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
			return fmt.Errorf("execute add books to collection: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("execute tx: %w", err)
	}

	return nil
}

// DeleteBooksCollection deletes books from a collection.
func (s *DB) DeleteBooksCollection(ctx context.Context, cID int64, bookIDs []int64) error {
	q := `DELETE FROM books_collection bc WHERE bc.collection_id = $1 AND bc.book_id = ANY ($2)`
	_, err := s.ExecContext(ctx, q, cID, pq.Array(bookIDs))
	if err != nil {
		return fmt.Errorf("remove books from collection: %w", err)
	}

	return nil
}
