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

	chunk := make([]byte, 0, chunkSize)
	parts := make([]storage.FilePart, fs.chunksCount)

	for {
		_, err := body.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		choosenStorage := fs.chooseStorage()
		path, err := choosenStorage.UploadFile(bytes.NewReader(chunk))
		parts = append(parts, storage.FilePart{Path: path, ContentLength: 1})
	}
	fs.metaStorage.Save(fileName, fileSize, parts)

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

func NewFileService(storages map[string]storage.Storage, metaStorage storage.MetaStorage, chunksCount int) IFileService {
	return &fileService{
		storages,
		metaStorage,
		chunksCount,
	}
}
