package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/caarlos0/env/v9"

	"github.com/Tsapen/bm/internal/bm"
)

type serverEnvs struct {
	RootDir string `env:"BM_ROOT_DIR"`
	Config  string `env:"BM_SERVER_CONFIG"`
}

type clientEnvs struct {
	RootDir string `env:"BM_ROOT_DIR"`
	Config  string `env:"BM_CLIENT_CONFIG"`
}

type ServerConfig struct {
	UnixSocketCfg
	HTTPCfg HTTPCfg `json:"http"`
	AraDB   DBCfg   `json:"db"`
}

type UnixSocketCfg struct {
	SocketPath   string `json:"socket_path"`
	ConnMaxCount int    `json:"connections_max_count"`

	Timeout time.Duration `json:"-"`
}

func (c *UnixSocketCfg) UnmarshalJSON(data []byte) error {
	cfg := new(struct {
		Timeout string `json:"timeout"`

		UnixSocketCfg
	})

	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	duration, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return err
	}

	c.Timeout = duration

	return nil
}

type HTTPCfg struct {
	Addr         string `json:"address"`
	ConnMaxCount int    `json:"connections_max_count"`

	Timeout time.Duration `json:"-"`
}

func (c *HTTPCfg) UnmarshalJSON(data []byte) error {
	cfg := new(struct {
		Timeout string `json:"timeout"`

		HTTPCfg
	})

	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	duration, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return err
	}

	c.Timeout = duration

	return nil
}

type DBCfg struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	HostName    string `json:"hostname"`
	Port        string `json:"port"`
	VirtualHost string `json:"virtual_host"`
}

type ClientConfig struct {
	Address string `json:"address"`

	Timeout time.Duration `json:"-"`
}

func (c *ClientConfig) UnmarshalJSON(data []byte) error {
	cfg := new(struct {
		Timeout string `json:"timeout"`
		ClientConfig
	})

	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return err
	}

	c.Timeout = timeout

	return nil
}

func GetForServer() (*ServerConfig, error) {
	envs := new(serverEnvs)
	if err := env.Parse(envs); err != nil {
		return nil, fmt.Errorf("get envs: %w", err)
	}

	cfg := new(ServerConfig)
	if err := readFromEnv(path.Join(envs.RootDir, envs.Config), cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}

func GetForClient() (*ClientConfig, error) {
	envs := new(clientEnvs)
	if err := env.Parse(envs); err != nil {
		return nil, fmt.Errorf("get envs: %w", err)
	}

	cfg := new(ClientConfig)
	if err := readFromEnv(path.Join(envs.RootDir, envs.Config), cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}

func readFromEnv(filepath string, receiver any) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("open file %s: %w", filepath, err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = bm.HandleErrPair(fmt.Errorf("close file: %w", closeErr), err)
		}
	}()

	if err = json.NewDecoder(file).Decode(receiver); err != nil {
		return fmt.Errorf("decode file: %w", err)
	}

	return
}
