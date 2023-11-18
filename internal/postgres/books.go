package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	bm "github.com/Tsapen/bm/internal/bm"
)

func (s *DB) CreateBook(ctx context.Context, b bm.Book) (int64, error) {
	query := `
		INSERT INTO books (title, author, published_date, edition, description, genre)
		VALUES (:title, :author, :published_date, :edition, :description, :genre)
		RETURNING id
	`

	var bookID int64
	err := s.QueryRowContext(ctx, query, b).Scan(&bookID)
	if err != nil {
		return 0, fmt.Errorf("insert book: %w", err)
	}

	return bookID, nil
}

func joinCollection(join bool) string {
	if !join {
		return ""
	}

	return `JOIN books_collection bc ON b.id=bc.book_id ` +
		`JOIN collections c ON bc.collection_id=c.id `
}

func booksWhereClause(f bm.BookFilter) (string, []any) {
	whereClauses := make([]string, 0)
	params := make([]any, 0)
	paramsCnt := 1
	if f.Author != "" {
		whereClauses = append(whereClauses, fmt.Sprint("b.author=$%d ", paramsCnt))
		params = append(params, f.Author)
		paramsCnt++
	}

	if f.Genre != "" {
		whereClauses = append(whereClauses, fmt.Sprint("b.genre=$%d ", paramsCnt))
		params = append(params, f.Genre)
		paramsCnt++
	}

	if !f.StartDate.IsZero() {
		whereClauses = append(whereClauses, fmt.Sprint("b.published_date>=$%d ", paramsCnt))
		params = append(params, f.StartDate)
		paramsCnt++
	}

	if !f.FinishDate.IsZero() {
		whereClauses = append(whereClauses, fmt.Sprint("b.published_date<=$%d ", paramsCnt))
		params = append(params, f.FinishDate)
		paramsCnt++
	}

	if len(whereClauses) == 0 {
		return "", nil
	}

	return fmt.Sprintf("WHERE (%s) ", strings.Join(whereClauses, "AND")), params
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
func (s *DB) Books(ctx context.Context, f bm.BookFilter) (books []bm.Book, err error) {
	q := "SELECT b.id, b.title, b.author, b.published_date, b.edition, b.description, b.genre FROM books b "
	q += joinCollection(f.CollectionID > 0)
	whereClause, params := booksWhereClause(f)
	q += whereClause
	q += orderBy("b", f.OrderBy, f.Desc)
	q += pagination(f.Page, f.PageSize)

	q, args, err := sqlx.Named(q, params)
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

	books = make([]bm.Book, 0, f.PageSize)
	for rows.Next() {
		var b bm.Book

		if err = rows.Scan(&b); err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (s *DB) UpdateBook(ctx context.Context, book bm.Book) error {
	result, err := s.DB.ExecContext(
		ctx, `
		UPDATE books SET
			author = :author,
			title = :title,
			edition = :edition,
			description = :description,
			published_date = :published_date,
			genre = :genre
		WHERE id = :id`,
		book,
	)
	if err != nil {
		return fmt.Errorf("update book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get the number of affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book with ID %d not found", book.ID)
	}

	return nil
}

func (s *DB) DeleteBooks(ctx context.Context, ids []int64) error {
	q := `DELETE FROM books b WHERE id IN ($1)`

	_, err := s.Exec(q, ids)
	return err
}
