package server

import (
	"filebalancer/internal/storages"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct {
	Server *http.Server
}

func CreateApp() *App {
	db, err := NewDb("filebalancer.db")
	if err != nil {
		panic(err)
	}
	SetupDb(db)
	// TODO check insert/get
	metaStorage := storages.FileStorage{Db: db}
	meta := storages.FileMetaArr{
		{
			ContentLength: 100,
			Number:        1,
			StorageUrl:    "test/123",
		},
	}
	file := &storages.File{
		Name:          "test",
		ContentLength: 10,
		FileMeta:      meta,
	}
	if err := metaStorage.Add(file); err != nil {
		panic(err)
	}
	res, err := metaStorage.Get(file.Name)
	if err != nil {
		panic(err)
	}
	log.Info().Msg(fmt.Sprintf("%v+", res))

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &App{
		Server: CreateHttpServer(),
	}
}
