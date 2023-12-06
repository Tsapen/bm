package httpclient

import (
	"context"
	"net"
	"net/http"
	"time"
)

// Config contains data for constructing client.
type Config struct {
	Address    string
	SocketPath string
	Timeout    time.Duration
}

// Clients communicates with BM http-server.
type Client struct {
	cfg Config

	httpClient *http.Client
}

// New constructs a new BM http-client.
func New(cfg Config) *Client {
	c := &http.Client{
		Timeout: cfg.Timeout,
	}

	if cfg.SocketPath != "" {
		c.Transport = &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", cfg.SocketPath)
			},
		}
	}

	return &Client{
		cfg:        cfg,
		httpClient: c,
	}
}
