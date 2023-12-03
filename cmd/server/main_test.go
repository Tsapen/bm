package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	bmtest "github.com/Tsapen/bm/cmd/server/bm-test"
	bmhttp "github.com/Tsapen/bm/internal/bm-http"
	bs "github.com/Tsapen/bm/internal/book-service"
	"github.com/Tsapen/bm/internal/config"
	"github.com/Tsapen/bm/internal/migrator"
	"github.com/Tsapen/bm/internal/postgres"
	"github.com/Tsapen/bm/pkg/api"
	httpclient "github.com/Tsapen/bm/pkg/http-client"
)

func TestBM(t *testing.T) {
	serverCfg, err := config.GetForServer()
	if err != nil {
		t.Fatalf("read server configs: %v\n", err)
	}

	clientCfg, err := config.GetForHTTPClient()
	if err != nil {
		t.Fatalf("read client configs: %v\n", err)
	}

	client := httpclient.New(httpclient.Config{
		Address: clientCfg.Address,
		Timeout: clientCfg.Timeout,
	})

	db, err := postgres.New(postgres.Config(*serverCfg.DB))
	if err != nil {
		t.Fatalf("init storage: %v\n", err)
	}

	cleanUpDB(t, serverCfg.DB)

	if err = migrator.New(serverCfg.MigrationsPath, db).Apply(); err != nil {
		log.Fatal().Err(err).Msg("apply migrations")
	}

	bookService := bs.New(db)

	httpService, err := bmhttp.NewServer(bmhttp.Config(*serverCfg.HTTPCfg), bookService)
	if err != nil {
		t.Fatalf("init http server: %v\n", err)
	}

	go func() {
		if err = httpService.StartTCPServer(); err != nil {
			t.Logf("run tcp server: %v\n", err)
		}
	}()

	waitRunning(t, client)

	bmtest.TestBM(t, client)
}

func waitRunning(t *testing.T, client *httpclient.Client) {
	const checkNum = 10
	const maxDelay = 100 * time.Millisecond

	ctx := context.Background()
	req := &api.GetBooksReq{}
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
		`DROP TABLE IF EXISTS books_collection`,
		`DROP TABLE IF EXISTS collections`,
		`DROP TABLE IF EXISTS books`,
	}
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			t.Fatalf("exec db query: %s", err)
		}
	}
}
