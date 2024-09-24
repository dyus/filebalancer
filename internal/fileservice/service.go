package fileservice

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/dyus/filebalancer/internal/balancer"
	"github.com/dyus/filebalancer/internal/db"
	"github.com/dyus/filebalancer/internal/storage"
	"github.com/google/uuid"
)

type IFileService interface {
	ReadFile(context.Context, string) (io.ReadCloser, error)
	SaveFile(context.Context, string, io.Reader, int64) error
}

type fileService struct {
	storageClient *storage.Client
	metaStorage   db.MetaStorage
	chunksCount   int64
	balancer      balancer.Balancer
}

func (fs *fileService) SaveFile(ctx context.Context, fileName string, body io.Reader, fileSize int64) error {
	chunkSizes, err := fs.calculateChunkSizes(fileSize)
	if err != nil {
		return err
	}

	hosts := fs.balancer.GetHosts(int(fs.chunksCount))
	chunks := make([]*db.FilePart, 0, len(chunkSizes))

	for i, size := range chunkSizes {
		part := &db.FilePart{
			Storage: hosts[i], Path: uuid.NewString(), ContentLength: size,
		}
		chunks = append(chunks, part)
	}

	if err = fs.metaStorage.Create(ctx, fileName, fileSize, chunks); err != nil {
		return err
	}

	for _, chunk := range chunks {
		limitBody := io.LimitReader(body, chunk.ContentLength)

		err := fs.storageClient.Write(ctx, chunk, limitBody)
		if err != nil {
			if cleanErr := fs.cleanup(ctx, fileName); cleanErr != nil {
				log.Printf("Can't cleanup data on bad save for file %s. error: %v\n", fileName, err)
			}

			return err
		}
	}

	if err = fs.metaStorage.Complete(ctx, fileName); err != nil {
		return err
	}

	return nil
}

func (fs *fileService) ReadFile(ctx context.Context, name string) (io.ReadCloser, error) {
	meta, err := fs.metaStorage.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	return &ChunksReader{
		meta.FileParts, fs.storageClient, nil, 0, ctx,
	}, nil
}

func (fs *fileService) calculateChunkSizes(fileSize int64) ([]int64, error) {
	if fs.chunksCount == 0 {
		log.Print("Wrong server configuration. ChanksCount can't equal 0")

		return nil, errors.New("chunks count can't equals 0")
	}

	chunkSize := fileSize / fs.chunksCount
	chunks := make([]int64, 0, fs.chunksCount)

	for i := range fs.chunksCount {
		if i == fs.chunksCount-1 && fileSize%chunkSize != 0 {
			chunkSize += fileSize % chunkSize
		}

		chunks = append(chunks, chunkSize)
	}

	return chunks, nil
}

func (fs *fileService) cleanup(ctx context.Context, fileName string) error {
	meta, err := fs.metaStorage.Get(ctx, fileName)
	if err != nil {
		return err
	}

	for _, chunk := range meta.FileParts {
		err = fs.storageClient.Delete(ctx, chunk)
		if err != nil {
			return err
		}
	}

	if err = fs.metaStorage.Error(ctx, meta.Name); err != nil {
		return err
	}

	return nil
}

type ChunksReader struct {
	parts         db.FilePartList
	storageClient *storage.Client
	currentFile   io.ReadCloser
	currentIndex  int
	ctx           context.Context
}

func (cr *ChunksReader) Read(data []byte) (int, error) {
	if cr.currentFile == nil {
		if err := cr.next(); err != nil {
			return 0, err
		}
	}

	n, err := cr.currentFile.Read(data)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return n, err
		}

		if err := cr.currentFile.Close(); err != nil {
			return n, err
		}

		if err := cr.next(); err != nil {
			return n, err
		}

		return n, nil
	}

	return n, nil
}

func (cr *ChunksReader) Close() error {
	if cr.currentFile != nil {
		return cr.currentFile.Close()
	}

	return nil
}

func (cr *ChunksReader) next() error {
	if cr.currentIndex >= len(cr.parts) {
		return io.EOF
	}

	part := cr.parts[cr.currentIndex]

	var err error

	cr.currentFile, err = cr.storageClient.Read(cr.ctx, part)
	if err != nil {
		return err
	}

	cr.currentIndex++

	return nil
}

func NewFileService(
	balancer balancer.Balancer,
	storageClient *storage.Client,
	metaStorage db.MetaStorage,
	chunksCount int64,
) IFileService {
	return &fileService{
		storageClient,
		metaStorage,
		chunksCount,
		balancer,
	}
}
