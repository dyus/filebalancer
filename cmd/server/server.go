package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dyus/filebalancer/app"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"golang.org/x/sync/errgroup"
)

func main() {
	configPath := "./configs/config.yaml"
	loader := confita.NewLoader(env.NewBackend(), file.NewBackend(configPath))
	conf := &app.Config{}

	if err := loader.Load(context.Background(), conf); err != nil {
		log.Fatalf("can't load config: %v\n", err)
	}

	srv, err := app.NewApplication(conf)
	if err != nil {
		log.Fatalf("can't run upp %v\n", err)
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case s := <-sig:
			log.Printf("got signal %s", s)

			return fmt.Errorf("signal stop: %s", s)
		}
	})
	g.Go(func() error {
		return srv.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
