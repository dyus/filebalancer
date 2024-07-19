package api

import (
	"net/http"

	"github.com/dyus/filebalancer/internal/file_service"
	"github.com/gorilla/mux"
)

func NewRouter(fileService file_service.IFileService) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/download/{path}", NewDownloadFileHandler(fileService)).Methods(http.MethodGet)
	r.HandleFunc("/upload/{path}", NewUploadFileHandler(fileService)).Methods(http.MethodPut)
	return r
}
