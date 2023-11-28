package cliclient

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/Tsapen/bm/internal/bm"
)

type (
	Config struct {
		SocketPath   string
		ConnMaxCount int64
		Timeout      time.Duration
	}

	CLIClient struct {
		cfg Config
	}
)

func New(cfg Config) *CLIClient {
	return &CLIClient{
		cfg: cfg,
	}
}

func (c *CLIClient) DoRequest(req any) (resp any, err error) {
	jsonCmd, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	conn, err := net.Dial("unix", c.cfg.SocketPath)
	if err != nil {
		return nil, fmt.Errorf("make connection: %w", err)
	}

	defer func() {
		err = bm.HandleErrPair(conn.Close(), err)
	}()

	if _, err = conn.Write(jsonCmd); err != nil {
		return nil, fmt.Errorf("write data: %w", err)
	}

	return resp, nil
}
