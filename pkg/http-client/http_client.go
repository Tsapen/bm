package httpclient

import (
	"net/http"
	"time"
)

type Config struct {
	Address string

	Timeout time.Duration
}

type Client struct {
	cfg Config

	httpClient *http.Client
}

func New(cfg Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}
