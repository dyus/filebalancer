package storage

import (
	"bytes"
	"io"
)

type InMemory struct {
	data map[string][]byte
}

func (i *InMemory) UploadFile(body io.Reader) (string, error) {
	path := "generate unique name"
	data, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	i.data[path] = data

	return path, nil
}

func (i *InMemory) ReadFile(path string) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(i.data[path]))
}

func NewInMemory() *InMemory {
	return &InMemory{
		data: map[string][]byte{},
	}
}
