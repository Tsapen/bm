package main

import (
	"context"
	"testing"
	"time"

	bmtest "github.com/Tsapen/bm/cmd/server/bm-test"
	"github.com/Tsapen/bm/internal/config"
	"github.com/Tsapen/bm/pkg/api"
	httpclient "github.com/Tsapen/bm/pkg/http-client"
)

func TestBM(t *testing.T) {
	clientCfg, err := config.GetForHTTPClient()
	if err != nil {
		t.Fatalf("read client configs: %v\n", err)
	}

	client := httpclient.New(httpclient.Config{
		Address: clientCfg.Address,
		Timeout: clientCfg.Timeout,
	})

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
