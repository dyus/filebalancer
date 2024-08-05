package file_service

import (
	"io"

	"github.com/dyus/filebalancer/internal/balancer"
	"github.com/dyus/filebalancer/internal/storage"
)

type IFileService interface {
	ReadFile(string) io.ReadCloser
	SaveFile(string, io.Reader, int64) error
}

type fileService struct {
	storageService *storage.StorageService
	metaStorage    storage.MetaStorage
	chunksCount    int64
	balancer       balancer.Balancer
}

func (fs *fileService) SaveFile(fileName string, body io.Reader, fileSize int64) error {
	chunkSize := fileSize / fs.chunksCount

	parts := make([]storage.FilePart, fs.chunksCount)
	hosts := fs.balancer.GetHosts(int(fs.chunksCount))
	storageIndex := 0
	for range fs.chunksCount {
		limit_body := io.LimitReader(body, chunkSize)
		choosenStorage := fs.storageService.Storages[hosts[storageIndex]]
		storageIndex += 1
		path, err := choosenStorage.UploadFile(limit_body)
		if err != nil {
			return err
		}
		parts = append(parts, storage.FilePart{Path: path, ContentLength: 1})
	}
	fs.metaStorage.Save(fileName, fileSize, parts)

	return nil
}

func (fs *fileService) ReadFile(name string) io.ReadCloser {
	// TODO get storages then compose file
	panic("not implemented")
}

func NewFileService(balancer balancer.Balancer, storageService *storage.StorageService, metaStorage storage.MetaStorage, chunksCount int64) IFileService {
	return &fileService{
		storageService,
		metaStorage,
		chunksCount,
		balancer,
	}
}
