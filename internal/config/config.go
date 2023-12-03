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
	RootDir        string `env:"BM_ROOT_DIR"`
	Config         string `env:"BM_SERVER_CONFIG"`
	MigrationsPath string `env:"BM_MIGRATIONS_PATH"`
}

type httpClientEnvs struct {
	RootDir string `env:"BM_ROOT_DIR"`
	Config  string `env:"BM_HTTP_CLIENT_CONFIG"`
}

type cliClientEnvs struct {
	RootDir string `env:"BM_ROOT_DIR"`
	Config  string `env:"BM_CLI_CLIENT_CONFIG"`
}

type ServerConfig struct {
	HTTPCfg *HTTPCfg `json:"http"`
	DB      *DBCfg   `json:"db"`

	MigrationsPath string `json:"-"`
}

type HTTPCfg struct {
	Addr         string `json:"address"`
	SocketPath   string `json:"socket_path"`
	ConnMaxCount int    `json:"connections_max_count"`

	Timeout time.Duration `json:"-"`
}

func (c *HTTPCfg) UnmarshalJSON(data []byte) error {
	type Alias HTTPCfg
	aux := &struct {
		Timeout string `json:"timeout"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	duration, err := time.ParseDuration(aux.Timeout)
	if err != nil {
		return fmt.Errorf("parse timeout: %w", err)
	}

	c.Timeout = duration

	return nil
}

type DBCfg struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	Port        string `json:"port"`
	VirtualHost string `json:"virtual_host"`

	HostName string `json:"host"`
}

type HTTPClientConfig struct {
	Address string `json:"address"`

	Timeout time.Duration `json:"-"`
}

func (c *HTTPClientConfig) UnmarshalJSON(data []byte) error {
	type Alias HTTPClientConfig
	aux := &struct {
		Timeout string `json:"timeout"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	duration, err := time.ParseDuration(aux.Timeout)
	if err != nil {
		return fmt.Errorf("parse timeout: %w", err)
	}

	c.Timeout = duration

	return nil
}

type CLIClientConfig struct {
	SocketPath string        `json:"socket_path"`
	Timeout    time.Duration `json:"-"`
}

func (c *CLIClientConfig) UnmarshalJSON(data []byte) error {
	type Alias CLIClientConfig
	aux := &struct {
		Timeout string `json:"timeout"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	duration, err := time.ParseDuration(aux.Timeout)
	if err != nil {
		return fmt.Errorf("parse timeout: %w", err)
	}

	c.Timeout = duration

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

	cfg.MigrationsPath = path.Join(envs.RootDir, envs.MigrationsPath)

	return cfg, nil
}

func GetForHTTPClient() (*HTTPClientConfig, error) {
	envs := new(httpClientEnvs)
	if err := env.Parse(envs); err != nil {
		return nil, fmt.Errorf("get envs: %w", err)
	}

	cfg := new(HTTPClientConfig)
	if err := readFromEnv(path.Join(envs.RootDir, envs.Config), cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}

func GetForCLIClient() (*CLIClientConfig, error) {
	envs := new(cliClientEnvs)
	if err := env.Parse(envs); err != nil {
		return nil, fmt.Errorf("get envs: %w", err)
	}

	cfg := new(CLIClientConfig)
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
