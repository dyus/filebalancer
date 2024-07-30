package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dyus/filebalancer/app"
	"github.com/dyus/filebalancer/internal/storage"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

func check_db(conf *app.Config) {
	db, err := app.NewDb(&conf.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	metaS := app.NewMetaStorage(db)
	parts := make([]storage.FilePart, 2)
	parts = append(parts,
		storage.FilePart{
			Path:          "test_path1",
			ContentLength: 20,
		},
	)
	parts = append(parts,
		storage.FilePart{
			Path:          "test_2",
			ContentLength: 20,
		},
	)
	err = metaS.Save(
		"test_name", 1024, parts,
	)
	if err != nil {
		log.Fatal(err)
	}

	file_meta, err := metaS.Get("test_name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", file_meta)
}

func main() {
	configPath := "./configs/config.yaml"
	loader := confita.NewLoader(env.NewBackend(), file.NewBackend(configPath))
	conf := &app.Config{}
	if err := loader.Load(context.Background(), conf); err != nil {
		fmt.Fprintf(os.Stderr, "can't load config: %s\n", err)
		os.Exit(1)
	}
	srv, err := app.NewApplication(conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server started at http://%s\n", conf.HTTP.Addr)
	log.Fatal(srv.Run())

	// check_db(conf)
}
