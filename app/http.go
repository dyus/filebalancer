package app

import (
	"net/http"

	"github.com/dyus/filebalancer/api"
	"github.com/dyus/filebalancer/internal/balancer"
	"github.com/dyus/filebalancer/internal/fileservice"
)

func NewHTTP(conf *HTTPConfig, fileService fileservice.IFileService, balancer balancer.Balancer) *http.Server {
	return &http.Server{
		Addr:    conf.Addr,
		Handler: api.NewRouter(fileService, balancer),
	}
}
