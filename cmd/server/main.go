package main

import (
	"github.com/rs/zerolog/log"

	bmhttp "github.com/Tsapen/bm/internal/bm-http"
	bmus "github.com/Tsapen/bm/internal/bm-unix-socket"
	bs "github.com/Tsapen/bm/internal/book-service"
	"github.com/Tsapen/bm/internal/config"
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

	bookService := bs.New(db)

	unixSocketServer, err := bmus.NewServer(bmus.Config(*cfg.UnixSocketCfg), bookService)
	if err != nil {
		log.Fatal().Err(err).Msg("init unix socket server")
	}

	go unixSocketServer.Start()

	httpService, err := bmhttp.NewServer(bmhttp.Config(*cfg.HTTPCfg), bookService)
	if err != nil {
		log.Fatal().Err(err).Msg("init http server")
	}

	if err = httpService.Start(); err != nil {
		log.Fatal().Err(err).Msg("run http server")
	}
}
