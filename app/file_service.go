package app

import (
	"github.com/dyus/filebalancer/internal/balancer"
	"github.com/dyus/filebalancer/internal/file_service"
	"github.com/dyus/filebalancer/internal/storage"
)

func newFileService(conf *Config, metaStorage storage.MetaStorage) file_service.IFileService {
	storages := make(map[string]storage.Storage)
	balancer := balancer.NewRoundRobinBalancer(conf.Balancer.Hosts)
	return file_service.NewFileService(balancer, storages, metaStorage, conf.ChunksCount)
}
