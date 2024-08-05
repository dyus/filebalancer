package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dyus/filebalancer/internal/file_service"
	"github.com/gorilla/mux"
)

func NewDownloadFileHandler(fileService file_service.IFileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", vars["path"]))
		file_reader := fileService.ReadFile(vars["path"])
		io.Copy(w, file_reader)
	}
}

func NewUploadFileHandler(fileService file_service.IFileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "blablabla"))
		vars := mux.Vars(r)
		fileService.SaveFile(vars["path"], r.Body, r.ContentLength)
		w.Write([]byte(vars["path"]))
	}
}
