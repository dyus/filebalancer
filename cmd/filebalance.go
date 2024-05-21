package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dyus/filebalancer/api"
)

func main() {
	srv := &http.Server{
		Handler:      api.NewRouter(),
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
