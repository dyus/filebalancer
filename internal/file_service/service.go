package file_service

import (
	"bytes"
	"io"

	"github.com/dyus/filebalancer/internal/storage"
)

type IFileService interface {
	ReadFile(string) io.ReadCloser
	SaveFile(string, io.Reader, int) error
}

type fileService struct {
	storages    map[string]storage.Storage
	metaStorage storage.MetaStorage
	chunksCount int
}

func (fs *fileService) SaveFile(fileName string, body io.Reader, fileSize int) error {
	chunkSize := fileSize / fs.chunksCount
	chunk := make([]byte, chunkSize)
	for {
		_, err := body.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fs.saveChunk(chunk, fs.chooseStorage())
	}
	return nil
}

func (fs *fileService) ReadFile(name string) io.ReadCloser {
	// TODO get storages then compose file
	panic("not implemented")
}

func (fs *fileService) chooseStorage() storage.Storage {
	// TODO choose somehow
	panic("not implemented")
}

func (fs *fileService) saveChunk(chunk []byte, storage storage.Storage) error {
	// TODO save meta
	// TODO save info about each saved object and clean on error
	path, err := storage.UploadFile(bytes.NewReader(chunk))
	if err != nil {
		return err
	}
	fs.metaStorage.Save(path)
	return nil
}

func NewFileService(storages map[string]storage.Storage, metaStorage storage.MetaStorage, chunksCount int) IFileService {
	return &fileService{
		storages,
		metaStorage,
		chunksCount,
	}
}
