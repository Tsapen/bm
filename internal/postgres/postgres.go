package postgres

import (
	"context"
	"database/sql"
	"fmt"

	bm "github.com/Tsapen/bm/internal/bm"
	_ "github.com/lib/pq"
)

// Config contains settings for db.
type Config struct {
	UserName    string
	Password    string
	HostName    string
	Port        string
	VirtualHost string
}

// DB contains db connection.
type DB struct {
	*sql.DB
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
	db, err := sql.Open("postgres", dbAddr)
	if err != nil {
		return nil, fmt.Errorf("can't open connection %s: %w", dbAddr, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't ping with connection %s: %w", dbAddr, err)
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
