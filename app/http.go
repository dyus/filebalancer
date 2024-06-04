package app

import (
	"net/http"

	"github.com/dyus/filebalancer/api"
	"github.com/dyus/filebalancer/internal/storage"
)

func NewHttp(conf *HTTPConfig, fileStorage *storage.InMemory) *http.Server {
	return &http.Server{
		Addr:    conf.Addr,
		Handler: api.NewRouter(fileStorage),
	}
}
