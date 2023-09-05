package storages

import (
	"bytes"
	"filebalancer/internal/services"
	"io"
	"sync"

	"github.com/google/uuid"
)

type Storages struct {
	hosts map[string]*InMemoryStorage
}

type InMemoryStorage struct {
	files map[string][]byte
	mutex sync.RWMutex
}

func (s *InMemoryStorage) Save(data io.Reader) (string, error) {
	// TODO move to variable
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

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{files: map[string][]byte{}}
}

func NewStorages() *Storages {
	return &Storages{hosts: map[string]*InMemoryStorage{}}
}

// TODO
// Make mechanic for choosing N storages based on total occupied bytes by each
func (s *Storages) GetStorages() []string {
	hosts := make([]string, 0, len(s.hosts))
	for k := range s.hosts {
		hosts = append(hosts, k)
	}
	return hosts
}

func (s *Storages) Register(host string) {
	s.hosts[host] = NewInMemoryStorage()
}

// TODO move to project structs ??
func (s *Storages) Save(parts []services.FilePart) {
	panic("not implemented")
}
