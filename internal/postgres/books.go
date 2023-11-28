package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	bm "github.com/Tsapen/bm/internal/bm"
)

func (s *DB) CreateBook(ctx context.Context, b bm.Book) (int64, error) {
	query := `
		INSERT INTO books (title, author, published_date, edition, description, genre)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var bookID int64
	err := s.QueryRowContext(
		ctx,
		query,
		b.Title,
		b.Author,
		b.PublishedDate,
		b.Edition,
		b.Description,
		b.Genre,
	).
		Scan(&bookID)
	if err != nil {
		return 0, fmt.Errorf("insert book: %w", err)
	}

	return bookID, nil
}

func joinCollection(f bm.BookFilter) string {
	if f.CollectionID == 0 {
		return ""
	}

	return `JOIN books_collection bc ON b.id=bc.book_id `
}

func booksWhereClause(f bm.BookFilter) (string, map[string]any) {
	whereClauses := make([]string, 0)
	params := make(map[string]any, 0)

	if f.ID != 0 {
		whereClauses = append(whereClauses, "b.id=:id ")
		params["id"] = f.ID
	}

	if f.Author != "" {
		whereClauses = append(whereClauses, "b.author=:author ")
		params["author"] = f.Author
	}

	if f.Genre != "" {
		whereClauses = append(whereClauses, "b.genre=:genre ")
		params["genre"] = f.Genre
	}

	if !f.StartDate.IsZero() {
		whereClauses = append(whereClauses, "b.published_date >= :start_date ")
		params["start_date"] = f.StartDate
	}

	if !f.FinishDate.IsZero() {
		whereClauses = append(whereClauses, "b.published_date <= :finish_date ")
		params["finish_date"] = f.FinishDate
	}

	if f.CollectionID != 0 {
		whereClauses = append(whereClauses, "bc.collection_id = :collection_id ")
		params["collection_id"] = f.CollectionID
	}

	if len(whereClauses) == 0 {
		return "", nil
	}

	return fmt.Sprintf("WHERE %s ", strings.Join(whereClauses, " AND ")), params
}

func orderBy(table, column string, desc bool) string {
	if column == "" {
		return ``
	}

	dir := ""
	if desc {
		dir = "DESC "
	}

	return fmt.Sprintf("ORDER BY %s.%s %s", table, column, dir)
}

func pagination(page, pageSize int64) string {
	return fmt.Sprintf("LIMIT %d OFFSET %d ", pageSize, (page-1)*pageSize)
}

// Books gets books by filter.
func (s *DB) Books(ctx context.Context, f bm.BookFilter) ([]bm.Book, error) {
	q := "SELECT b.id, b.title, b.author, b.published_date, b.edition, b.description, b.genre FROM books b "
	q += joinCollection(f)
	whereClause, params := booksWhereClause(f)
	q += whereClause
	q += orderBy("b", f.OrderBy, f.Desc)
	q += pagination(f.Page, f.PageSize)

	rows, err := s.NamedQueryContext(ctx, q, params)
	if err != nil {
		return nil, fmt.Errorf("select books: %w", err)
	}

	defer func() {
		err = bm.HandleErrPair(rows.Close(), err)
	}()

	var books []bm.Book
	if err = sqlx.StructScan(rows, &books); err != nil {
		return nil, fmt.Errorf("copy data into struct: %w", err)
	}

	return books, nil
}

func (s *DB) UpdateBook(ctx context.Context, b bm.Book) error {
	params := []any{b.Author, b.Title, b.Edition, b.Description, b.PublishedDate, b.Genre, b.ID}
	q := `UPDATE books SET
			author = $1,
			title = $2,
			edition = $3,
			description = $4,
			published_date = $5,
			genre = $6
		WHERE id = $7`

	result, err := s.DB.ExecContext(ctx, q, params...)
	if err != nil {
		return fmt.Errorf("update book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get the number of affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book with ID %d not found", b.ID)
	}

	return nil
}

func (s *DB) DeleteBooks(ctx context.Context, ids []int64) error {
	q := `DELETE FROM books b WHERE id = ANY($1)`

	_, err := s.ExecContext(ctx, q, pq.Array(ids))
	return err
}
