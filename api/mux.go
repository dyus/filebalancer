package api

import (
	"net/http"

	"github.com/dyus/filebalancer/internal/balancer"
	"github.com/dyus/filebalancer/internal/fileservice"
	"github.com/gorilla/mux"
)

func NewRouter(fileService fileservice.IFileService, balancer balancer.Balancer) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/download/{path}", DownloadFileHandler(fileService)).Methods(http.MethodGet)
	r.HandleFunc("/upload/{path}", UploadFileHandler(fileService)).Methods(http.MethodPut)
	r.HandleFunc("/storage/register", StorageRegisterHandler(balancer)).Methods(http.MethodPut)

	return r
}
