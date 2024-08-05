package app

import (
	"github.com/dyus/filebalancer/internal/balancer"
	"github.com/dyus/filebalancer/internal/file_service"
	"github.com/dyus/filebalancer/internal/storage"
)

func newFileService(conf *Config, metaStorage storage.MetaStorage) (file_service.IFileService, error) {
	storageService, err := storage.NewStorageService("InMemory", conf.Balancer.Hosts)
	if err != nil {
		return nil, err
	}
	balancer := balancer.NewRoundRobinBalancer(conf.Balancer.Hosts)
	return file_service.NewFileService(balancer, storageService, metaStorage, conf.ChunksCount), nil
}
