package app

import (
	"net/http"

	"github.com/dyus/filebalancer/api"
)

func NewHttp(conf *HTTPConfig) *http.Server {
	return &http.Server{
		Addr:    conf.Addr,
		Handler: api.NewRouter(),
	}
}
