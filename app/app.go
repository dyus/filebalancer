package app

import (
	"net/http"
	"time"
)

type Application struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func (m *Application) Run() error {
	return m.server.ListenAndServe()
}
func NewApplication(conf *Config) (*Application, error) {
	db, err := newDb(&conf.DBConfig)
	if err != nil {
		return nil, err
	}
	metaStorage := newMetaStorage(db)
	fileService := newFileService(conf, metaStorage)
	return &Application{
		server:          NewHttp(&conf.HTTP, fileService),
		shutdownTimeout: conf.ShutdownTimeout,
	}, nil
}
