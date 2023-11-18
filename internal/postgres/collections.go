package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	bm "github.com/Tsapen/bm/internal/bm"
)

// Collections gets collections.
func (s *DB) Collections(ctx context.Context, f bm.CollectionsFilter) ([]bm.Collection, error) {
	q := "SELECT c.id, c.name, c.description FROM c collections "
	q += collectionsWhereClause(f)
	q += orderBy("c", f.OrderBy, f.Desc)
	q += pagination(f.Page, f.PageSize)

	q, args, err := sqlx.Named(q, f)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = s.Query(q, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = bm.HandleErrPair(rows.Close(), err)
	}()

	collections := make([]bm.Collection, 0, f.PageSize)
	for rows.Next() {
		var c bm.Collection
		if err = rows.Scan(&c); err != nil {
			return nil, err
		}

		collections = append(collections, c)
	}

	return collections, nil
}

func collectionsWhereClause(f bm.CollectionsFilter) string {
	whereClauses := make([]string, 0)

	if len(f.IDs) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("id IN ($1)", f.IDs))
	}

	if len(whereClauses) == 0 {
		return ""
	}

	return "WHERE " + strings.Join(whereClauses, " AND ")
}

func (s *DB) CreateCollection(ctx context.Context, c bm.Collection) (int64, error) {
	query := `
		INSERT INTO collections (name, description)
		VALUES (:name, :description)
		RETURNING id
	`

	var id int64
	err := s.QueryRowContext(ctx, query, c).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("execute collection insert: %w", err)
	}

	return id, nil
}

func insertBooksCollectionValues(numValues int) string {
	values := make([]string, 0, numValues)
	for i := 0; i < numValues; i++ {
		values = append(values, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
	}

	return strings.Join(values, ", ")
}

func collectionFieldsForUpdate(c bm.Collection) string {
	q := []string{}
	if c.Name != "" {
		q = append(q, "name = :name")
	}

	if c.Description != "" {
		q = append(q, "description = :description")
	}

	return strings.Join(q, ", ")
}

// UpdateCollection updates a collection and its books.
func (s *DB) UpdateCollection(ctx context.Context, c bm.Collection) (err error) {
	q := `UPDATE collections SET ` + collectionFieldsForUpdate(c) + ` WHERE id = :id`
	err = s.QueryRowContext(ctx, q, c).Err()
	if err != nil {
		return fmt.Errorf("update collection: %w", err)
	}

	return nil
}

// DeleteCollections deletes collection and its associations.
func (s *DB) DeleteCollection(ctx context.Context, id int64) error {
	err := s.withTX(ctx, func(tx *sql.Tx) error {
		q := `DELETE FROM collection_books WHERE collection_id = $1`
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

// CreateBooksCollection adds books to a collection.
func (s *DB) CreateBooksCollection(ctx context.Context, cID int64, bookIDs []int64) error {
	err := s.withTX(ctx, func(tx *sql.Tx) error {
		q := fmt.Sprintf(`
		INSERT INTO collection_books (collection_id, book_id)
		VALUES %s
	`, insertBooksCollectionValues(len(bookIDs)))

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
	q := `DELETE FROM collection_books WHERE collection_id = $1 AND book_id IN ($2)`
	_, err := s.ExecContext(ctx, q, cID, bookIDs)
	if err != nil {
		return fmt.Errorf("remove books from collection: %w", err)
	}

	return nil
}
