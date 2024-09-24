package storage

import (
	"bytes"
	"fmt"
	"io"
	"sync/atomic"
)

type InMemory struct {
	usedSpace int64
	data      map[string][]byte
}

func (i *InMemory) UploadFile(path string, body io.Reader) (string, int64, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return "", 0, err
	}

	i.data[path] = data
	fileSize := int64(len(data))
	i.updateUsedSpace(int(fileSize))

	return path, fileSize, nil
}

func (i *InMemory) ReadFile(path string) (io.ReadCloser, error) {
	if _, exists := i.data[path]; !exists {
		return nil, fmt.Errorf("file with name %s", path)
	}

	return io.NopCloser(bytes.NewReader(i.data[path])), nil
}

func (i *InMemory) DeleteFile(path string) {
	if data, exists := i.data[path]; exists {
		i.updateUsedSpace(-len(data))
		delete(i.data, path)
	}
}

func (i *InMemory) Status() *Status {
	return &Status{
		UsedSpace: i.usedSpace,
	}
}

func (i *InMemory) updateUsedSpace(fileSize int) {
	atomic.AddInt64(&i.usedSpace, int64(fileSize))
}

func NewInMemory() *InMemory {
	return &InMemory{
		data: map[string][]byte{},
	}
}
