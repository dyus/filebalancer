package app

import (
	"net/http"

	"github.com/dyus/filebalancer/api"
	"github.com/dyus/filebalancer/internal/file_service"
)

func NewHttp(conf *HTTPConfig, fileService file_service.IFileService) *http.Server {
	return &http.Server{
		Addr:    conf.Addr,
		Handler: api.NewRouter(fileService),
	}
}
