package storage

import (
	"bytes"
	"io"
)

type InMemory struct {
	data map[string][]byte
}

func (i *InMemory) UploadFile(path string, body io.Reader) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	i.data[path] = data

	return nil
}

func (i *InMemory) ReadFile(path string) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(i.data[path]))
}

func NewInMemory() *InMemory {
	return &InMemory{
		data: map[string][]byte{},
	}
}
