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
	db, err := NewDb(&conf.DBConfig)
	if err != nil {
		return nil, err
	}
	metaStorage := NewMetaStorage(db)
	fileService, err := newFileService(conf, metaStorage)
	if err != nil {
		return nil, err
	}
	return &Application{
		server:          NewHttp(&conf.HTTP, fileService),
		shutdownTimeout: conf.ShutdownTimeout,
	}, nil
}
