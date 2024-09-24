package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dyus/filebalancer/internal/storage"
	"github.com/gorilla/mux"
)

type StorageResponse struct {
	Path   string `json:"path"`
	Length int64  `json:"length"`
}

func saveHandler(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		path, length, err := storage.UploadFile(vars["path"], r.Body)
		if err != nil {
			log.Printf("Can't save file %v", err)
			http.Error(w, "Can't save file", http.StatusInternalServerError)

			return
		}

		s := StorageResponse{
			path,
			length,
		}

		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Printf("Can't encode response %v", err)
			http.Error(w, "Can't encode response", http.StatusInternalServerError)

			return
		}

		_, err = w.Write([]byte(path))
		if err != nil {
			log.Printf("Can't write response %v", err)
			http.Error(w, fmt.Sprintf("Can't write response %v", err), http.StatusInternalServerError)

			return
		}
	}
}

func getHandler(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		reader, err := storage.ReadFile(vars["path"])
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)

			return
		}

		length, err := io.Copy(w, reader)
		if err != nil {
			http.Error(w, "can't write file", http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Length", strconv.FormatInt(length, 10))
	}
}

func statusHandler(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(storage.Status()); err != nil {
			log.Printf("Error in statusHandler. Can't encode json %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}
}

func main() {
	router := mux.NewRouter()
	storage := storage.NewInMemory()
	router.HandleFunc("/download/{path}", getHandler(storage)).Methods(http.MethodGet)
	router.HandleFunc("/upload/{path}", saveHandler(storage)).Methods(http.MethodPut)
	router.HandleFunc("/status", statusHandler(storage)).Methods(http.MethodGet)

	addr := os.Getenv("STORAGE_ADDR")
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("Server started at http://%s\n", addr)
	log.Fatal(server.ListenAndServe())
}
