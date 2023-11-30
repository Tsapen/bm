package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	bm "github.com/Tsapen/bm/internal/bm"
)

// Config contains settings for db.
type Config struct {
	UserName    string
	Password    string
	Port        string
	VirtualHost string
	HostName    string
}

// DB contains db connection.
type DB struct {
	*sqlx.DB
}

func (c *Config) dbAddr() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		c.UserName,
		c.Password,
		c.HostName,
		c.Port,
		c.VirtualHost,
	)
}

// New create new storage.
func New(c Config) (*DB, error) {
	dbAddr := c.dbAddr()
	db, err := sqlx.Open("postgres", dbAddr)
	if err != nil {
		return nil, fmt.Errorf("open connection %s: %w", dbAddr, err)
	}

	for i := 0; i < 10; i++ {
		if err = db.Ping(); err == nil {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	if err != nil {
		return nil, fmt.Errorf("ping with connection %s: %w", dbAddr, err)
	}

	return &DB{
		db,
	}, nil
}

func (s *DB) withTX(ctx context.Context, fnc func(tx *sql.Tx) error) error {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}

	defer func() {
		if err != nil {
			err = bm.HandleErrPair(tx.Rollback(), err)
		} else {
			err = bm.HandleErrPair(tx.Commit(), err)
		}
	}()

	return fnc(tx)
}
