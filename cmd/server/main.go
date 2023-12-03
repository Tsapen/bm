package main

import (
	"github.com/rs/zerolog/log"

	bmhttp "github.com/Tsapen/bm/internal/bm-http"
	bs "github.com/Tsapen/bm/internal/book-service"
	"github.com/Tsapen/bm/internal/config"
	"github.com/Tsapen/bm/internal/migrator"
	"github.com/Tsapen/bm/internal/postgres"
)

func main() {
	cfg, err := config.GetForServer()
	if err != nil {
		log.Fatal().Err(err).Msg("read config")
	}

	db, err := postgres.New(postgres.Config(*cfg.DB))
	if err != nil {
		log.Fatal().Err(err).Msg("init storage")
	}

	if err = migrator.New(cfg.MigrationsPath, db).Apply(); err != nil {
		log.Fatal().Err(err).Msg("apply migrations")
	}

	bookService := bs.New(db)

	httpService, err := bmhttp.NewServer(bmhttp.Config(*cfg.HTTPCfg), bookService)
	if err != nil {
		log.Fatal().Err(err).Msg("init http server")
	}

	go func() {
		if err = httpService.StartUnixSocketServer(); err != nil {
			log.Fatal().Err(err).Msg("run unix socket server server")
		}
	}()

	if err = httpService.StartTCPServer(); err != nil {
		log.Fatal().Err(err).Msg("run tcp server")
	}
}
