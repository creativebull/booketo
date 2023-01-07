package app

import (
	"os"

	"github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func StartApp() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	runGinServer(config)
}

func runGinServer(config utils.Config) {
	server, err := NewServer(config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
