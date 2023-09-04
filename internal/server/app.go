package server

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct {
	Server *http.Server
}

func CreateApp() *App {
	db, err := NewDb()
	if err != nil {
		panic(err)
	}
	SetupDb(db)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &App{
		Server: CreateHttpServer(),
	}
}
