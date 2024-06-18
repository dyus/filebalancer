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
	storages map[int]Storage
}

func NewStorageService(storageType string, num int) (*StorageService, error) {
	if storageType != "InMemory" {
		return nil, errors.New(fmt.Sprintf("Unknown type %s", storageType))
	}

	storages := make(map[int]Storage)
	for num := range num {
		storages[num] = NewInMemory()
	}

	return &StorageService{
		storages: storages,
	}, nil
}
