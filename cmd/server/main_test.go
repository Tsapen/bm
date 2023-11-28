package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	bmtest "github.com/Tsapen/bm/cmd/server/bm-test"
	bmhttp "github.com/Tsapen/bm/internal/bm-http"
	bs "github.com/Tsapen/bm/internal/book-service"
	"github.com/Tsapen/bm/internal/config"
	"github.com/Tsapen/bm/internal/postgres"
	"github.com/Tsapen/bm/pkg/api"
	httpclient "github.com/Tsapen/bm/pkg/http-client"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestBM(t *testing.T) {
	serverCfg, err := config.GetForServer()
	if err != nil {
		t.Fatalf("read server configs: %v\n", err)
	}

	clientCfg, err := config.GetForHTTPClient()
	if err != nil {
		t.Fatalf("read client configs: %v\n", err)
	}

	client := httpclient.New(httpclient.Config(*clientCfg))

	db, err := postgres.New(postgres.Config(*serverCfg.DB))
	if err != nil {
		t.Fatalf("init storage: %v\n", err)
	}

	bookService := bs.New(db)

	httpService, err := bmhttp.NewServer(bmhttp.Config(*serverCfg.HTTPCfg), bookService)
	if err != nil {
		t.Fatalf("init http server: %v\n", err)
	}

	cleanUpDB(t, serverCfg.DB)

	go func() {
		if err = httpService.Start(); err != nil {
			t.Logf("run http server: %v\n", err)
		}
	}()

	waitRunning(t, client)

	bmtest.TestBM(t, client)
}

func waitRunning(t *testing.T, client *httpclient.Client) {
	const checkNum = 10
	const maxDelay = 100 * time.Millisecond

	ctx := context.Background()
	req := &api.GetBooksReq{
		Genre:    "horror",
		Page:     1,
		PageSize: 1,
	}
	for i := 0; i < checkNum; i++ {
		time.Sleep(maxDelay)

		if _, err := client.GetBooks(ctx, req); err == nil {
			return
		}
	}

	t.Fatalf("service could not start")
}

func dbAddr(c *config.DBCfg) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		c.UserName,
		c.Password,
		c.HostName,
		c.Port,
		c.VirtualHost,
	)
}

func cleanUpDB(t *testing.T, c *config.DBCfg) {
	db, err := sql.Open("postgres", dbAddr(c))
	if err != nil {
		t.Fatalf("open connection: %s", err)
	}

	var queries = []string{
		`DROP TABLE IF EXISTS migrations`,
		`DROP TABLE IF EXISTS books_collection`,
		`DROP TABLE IF EXISTS collections`,
		`DROP TABLE IF EXISTS books`,
		`CREATE TABLE IF NOT EXISTS books (
			id SERIAL NOT NULL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			author VARCHAR(100) NOT NULL,
			genre VARCHAR(100) NOT NULL,
			published_date timestamp,
			edition VARCHAR(100) NOT NULL,
			description TEXT
		);
		
		CREATE TABLE IF NOT EXISTS collections (
			id SERIAL NOT NULL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT
		);
		
		CREATE TABLE IF NOT EXISTS books_collection (
			collection_id INT NOT NULL REFERENCES collections(id),
			book_id INT NOT NULL REFERENCES books(id),
			PRIMARY KEY(book_id, collection_id)
		);
		`,
	}
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			t.Fatalf("exec db query: %s", err)
		}
	}
}
