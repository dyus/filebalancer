package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/download", NewDownloadFileHandler).Methods(http.MethodGet)
	r.HandleFunc("/upload", NewUploadFileHandler).Methods(http.MethodPut)
	return r
}
