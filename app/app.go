package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dyus/filebalancer/api"
	"github.com/dyus/filebalancer/internal/balancer"
	"github.com/dyus/filebalancer/internal/db"
	"github.com/dyus/filebalancer/internal/fileservice"
	"github.com/dyus/filebalancer/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // pg driver
	"golang.org/x/sync/errgroup"
)

const CleanUpTime = 24 * time.Hour

type Application struct {
	server          *http.Server
	fileService     fileservice.IFileService
	shutdownTimeout time.Duration
}

func (m *Application) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Printf("Server started at http://%s\n", m.server.Addr)

		if err := m.server.ListenAndServe(); err != nil {
			return fmt.Errorf("listen and serve error: %w", err)
		}

		return nil
	})

	cleanUpTicker := time.NewTicker(CleanUpTime)
	defer cleanUpTicker.Stop()

	g.Go(func() error {
		err := m.fileService.Cleanup(ctx, cleanUpTicker)
		if err != nil {
			return fmt.Errorf("can't execute cleanup task: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		log.Printf("Graceful shutdown")

		shutDownCtx, cancel := context.WithTimeout(ctx, m.shutdownTimeout)

		defer cancel()

		if err := m.server.Shutdown(shutDownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}

		return ctx.Err()
	})

	return g.Wait()
}

func newDB(conf *DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			conf.Username, conf.Password, conf.Name, conf.Host, conf.Port))
	if err != nil {
		return nil, err
	}

	sql, err := os.ReadFile("db/create_db.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewApplication(conf *Config) (*Application, error) {
	dbConnect, err := newDB(&conf.DBConfig)
	if err != nil {
		return nil, err
	}

	metaStorage := db.NewMetaStorage(dbConnect)
	balancer := balancer.NewRoundRobinBalancer(conf.Balancer.Hosts)
	storageClient := storage.NewClient()
	fileService := fileservice.NewFileService(balancer, storageClient, metaStorage, conf.ChunksCount)

	return &Application{
		server: &http.Server{
			Addr:    conf.HTTP.Addr,
			Handler: api.NewRouter(fileService, balancer),
		},
		fileService:     fileService,
		shutdownTimeout: conf.ShutdownTimeout,
	}, nil
}
