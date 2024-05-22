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
	config := &app.Config{}
	if err := loader.Load(context.Background(), config); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "can't load config: %s\n", err)
		os.Exit(1)
	}

	srv := app.NewHttp(&config.HTTP)
	log.Fatal(srv.ListenAndServe())
}
