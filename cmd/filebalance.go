package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dyus/filebalancer/app"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

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
}
