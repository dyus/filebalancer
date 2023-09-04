package storages

import (
	"bytes"
	"io"
	"sync"

	"github.com/google/uuid"
)

type InMemoryStorage struct {
	files map[string][]byte
	mutex sync.RWMutex
}

func (s *InMemoryStorage) Save(data io.Reader) (string, error) {
	fileName := uuid.NewString()[:10]
	fileBytes, err := io.ReadAll(data)
	if err != nil {
		return "", err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.files[fileName] = fileBytes

	return fileName, nil
}

func (s *InMemoryStorage) Read(fileName string) io.ReadCloser {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return io.NopCloser(bytes.NewReader(s.files[fileName]))
}

func New() *InMemoryStorage {
	return &InMemoryStorage{files: map[string][]byte{}}
}
