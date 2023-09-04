package server

import (
	"filebalancer/internal/middleware"
	"filebalancer/internal/storages"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type FileMeta struct {
	Path string
	Name string
}

func CreateHttpServer() *http.Server {
	
	log.Info().Msg("start server")

	router := mux.NewRouter()
	storage := storages.New()

	router.HandleFunc("/upload/{name}", uploadFileHandler(storage))
	router.HandleFunc("/download/{name}", downloadFileHandler(storage))

	router.Use(middleware.Logging)

	http.Handle("/", router)

	return &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}

func uploadFileHandler(storage *storages.InMemoryStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filePath, err := storage.Save(r.Body)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		fileMeta := FileMeta{
			Path: filePath,
			Name: vars["name"],
		}
		io.WriteString(w, fmt.Sprintf("%+v\n", fileMeta))
	}
}

func downloadFileHandler(storage *storages.InMemoryStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		io.Copy(w, storage.Read(vars["name"]))
	}
}
