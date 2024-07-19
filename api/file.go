package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

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
		length, err := strconv.Atoi(r.Header.Get("Content-Length"))
		if err != nil {
			http.Error(w, "Can't get content length", http.StatusInternalServerError)
		}
		fileService.SaveFile(vars["path"], r.Body, length)
		w.Write([]byte(vars["path"]))
	}
}
