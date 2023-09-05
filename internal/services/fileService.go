package services

import (
	"filebalancer/internal/storages"
	"fmt"
	"io"
	"math"

	"github.com/google/uuid"
)

type FilePart struct {
	ContentLength int
	Number        int
	StorageUrl    string
}

type FileData struct {
	name          string
	ContentLength int
	reader        io.ByteReader
}

type FileService struct {
	chunksNum    int
	fileStorages storages.Storages
	metaStorage  storages.FileStorage
}

func (fs *FileService) Save(fileData *FileData) error {
	// save FileParts and file data to db
	// save files parts to storage
	parts := fs.makeChunks(fileData)
	error := fs.metaStorage.Add(parts)
	if error != nil {
		return error
	}
	fs.fileStorages.Save(parts)

}

func (fs *FileService) makeChunks(fileData *FileData) []FilePart {
	chunkSize := fs.getChunkSize(fileData.ContentLength)
	parts := make([]FilePart, 0, fs.chunksNum)
	choosenStorages := fs.fileStorages.GetStorages()
	lastBytes := fileData.ContentLength
	index := 0
	// TODO get available storages throw balancer (write it)
	for lastBytes > 0 {
		if lastBytes > chunkSize {
			parts = append(parts, FilePart{
				ContentLength: chunkSize,
				Number:        index,
				StorageUrl:    fmt.Sprintf("%s/%s", choosenStorages[index], uuid.NewString()[:10]),
			})
		} else {
			parts = append(parts, FilePart{
				ContentLength: lastBytes,
				Number:        index,
				StorageUrl:    fmt.Sprintf("%s/%s", choosenStorages[index], uuid.NewString()[:10]),
			})
		}
		lastBytes -= chunkSize
		index += 1
	}
	return parts
}

func (fs *FileService) getChunkSize(contentLength int) int {
	return int(math.Ceil(float64(contentLength) / float64(fs.chunksNum)))
}

func NewFileService() *FileService {
	return &FileService{
		chunksNum: 5,
	}
}
