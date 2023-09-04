package main

import (
	"filebalancer/internal/server"

	"github.com/rs/zerolog/log"
)

func main() {
	err := server.CreateApp().Server.ListenAndServe()
	if err != nil {
		log.Error().Err(err).Msg("stop server")
	} else {
		log.Info().Msg("stop server")
	}
}
