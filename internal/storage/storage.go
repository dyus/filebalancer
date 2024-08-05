package storage

import (
	"errors"
	"fmt"
	"io"
)

type Storage interface {
	UploadFile(io.Reader) (string, error)
	ReadFile(string) io.ReadCloser
}

type StorageService struct {
	Storages map[string]Storage
}

func NewStorageService(storageType string, names []string) (*StorageService, error) {
	if storageType != "InMemory" {
		return nil, errors.New(fmt.Sprintf("Unknown type %s", storageType))
	}

	storages := make(map[string]Storage)
	for _, name := range names {
		storages[name] = NewInMemory()
	}

	return &StorageService{
		Storages: storages,
	}, nil
}
