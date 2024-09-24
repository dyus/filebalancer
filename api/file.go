package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/dyus/filebalancer/internal/balancer"
	"github.com/dyus/filebalancer/internal/fileservice"
	"github.com/gorilla/mux"
)

func DownloadFileHandler(fileService fileservice.IFileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Header().Set("Content-Disposition", "attachment; filename="+vars["path"])

		fileReader, err := fileService.ReadFile(r.Context(), vars["path"])
		if err != nil {
			log.Printf("Can't download file %v", err)
			http.Error(w, fmt.Sprintf("Can't download file %v", err), http.StatusInternalServerError)

			return
		}

		_, err = io.Copy(w, fileReader)
		if err != nil {
			log.Printf("Can't write response %v", err)
			http.Error(w, fmt.Sprintf("Can't write response %v", err), http.StatusInternalServerError)

			return
		}
	}
}

func UploadFileHandler(fileService fileservice.IFileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		err := fileService.SaveFile(r.Context(), vars["path"], r.Body, r.ContentLength)
		if err != nil {
			log.Printf("Can't upload file %v", err)
			http.Error(w, fmt.Sprintf("Can't upload file %v", err), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+vars["path"])

		_, err = w.Write([]byte(vars["path"]))
		if err != nil {
			log.Printf("Can't write response %v", err)
			http.Error(w, fmt.Sprintf("Can't write response %v", err), http.StatusInternalServerError)

			return
		}
	}
}

type StorageRegister struct {
	Host string `json:"host"`
}

func StorageRegisterHandler(balancer balancer.Balancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		var sr StorageRegister

		err := decoder.Decode(&sr)
		if err != nil {
			log.Printf("Can't add storage %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)

			return
		}

		queryURL, err := url.JoinPath(sr.Host, "status")
		if err != nil {
			log.Printf("Can't add storage %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)

			return
		}

		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, queryURL, nil)
		if err != nil {
			log.Printf("Can't add storage %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)

			return
		}

		resp, err := http.DefaultClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}

		if err != nil {
			log.Printf("Can't add storage %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)

			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Can't add storage %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)

			return
		}

		balancer.AddHost(sr.Host)
	}
}
