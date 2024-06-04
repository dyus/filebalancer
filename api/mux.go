package api

import (
	"net/http"

	"github.com/dyus/filebalancer/internal/storage"
	"github.com/gorilla/mux"
)

func NewRouter(fileStorage *storage.InMemory) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/download/{path}", NewDownloadFileHandler(fileStorage)).Methods(http.MethodGet)
	r.HandleFunc("/upload/{path}", NewUploadFileHandler(fileStorage)).Methods(http.MethodPut)
	return r
}
