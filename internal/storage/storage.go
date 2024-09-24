package storage

import "io"

type Storage interface {
	UploadFile(string, io.Reader) (string, int64, error)
	ReadFile(string) (io.ReadCloser, error)
	DeleteFile(string)
	Status() *Status
}
