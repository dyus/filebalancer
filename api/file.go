package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dyus/filebalancer/internal/storage"
	"github.com/gorilla/mux"
)

func NewDownloadFileHandler(storage *storage.InMemory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", vars["path"]))
		file_reader := storage.ReadFile(vars["path"])
		io.Copy(w, file_reader)
	}
}

func NewUploadFileHandler(storage *storage.InMemory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "blablabla"))
		vars := mux.Vars(r)
		storage.UploadFile(vars["path"], r.Body)
		w.Write([]byte(vars["path"]))
	}
}
