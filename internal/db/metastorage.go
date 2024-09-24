package db

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FileStatus int

const (
	InProgress FileStatus = iota
	Complete
	Error
)

type MetaStorage interface {
	Create(context.Context, string, int64, []*FilePart) error
	Complete(context.Context, string) error
	Error(context.Context, string) error
	Get(context.Context, string) (*fileMeta, error)
}

type pgMetaStorage struct {
	db *sqlx.DB
}

type FilePart struct {
	Storage       string `db:"storage"`
	Path          string `db:"path"`
	ContentLength int64  `db:"content_length"`
}

type FilePartList []*FilePart

func (f FilePart) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FilePart) Scan(val interface{}) error {
	data, ok := val.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", val)
	}

	return json.Unmarshal(data, &f)
}

func (f FilePartList) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FilePartList) Scan(val interface{}) error {
	data, ok := val.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", val)
	}

	return json.Unmarshal(data, &f)
}

type fileMeta struct {
	ID            int          `db:"id"`
	Name          string       `db:"name"`
	ContentLength int64        `db:"content_length"`
	FileParts     FilePartList `db:"parts"`
	Status        FileStatus   `db:"status"`
}

func (m *pgMetaStorage) Create(ctx context.Context, path string, contentLength int64, parts []*FilePart) error {
	fileMeta := fileMeta{Name: path, ContentLength: contentLength, FileParts: parts}
	_, err := m.db.ExecContext(ctx, `INSERT INTO file_meta (name, content_length, parts, status)
		VALUES ($1, $2, $3, $4)`,
		fileMeta.Name, fileMeta.ContentLength, fileMeta.FileParts, InProgress)

	return err
}

func (m *pgMetaStorage) Complete(ctx context.Context, path string) error {
	_, err := m.db.ExecContext(ctx, `UPDATE file_meta SET status=$1 WHERE name=$2`, Complete, path)

	return err
}

func (m *pgMetaStorage) Error(ctx context.Context, path string) error {
	_, err := m.db.ExecContext(ctx, `UPDATE file_meta SET status=$1 WHERE name=$2`, Error, path)

	return err
}

func (m *pgMetaStorage) Get(ctx context.Context, name string) (*fileMeta, error) {
	fileMeta := fileMeta{}

	err := m.db.GetContext(ctx, &fileMeta, "SELECT * FROM file_meta WHERE name=$1", name)
	if err != nil {
		return nil, err
	}

	return &fileMeta, nil
}

func NewMetaStorage(db *sqlx.DB) MetaStorage {
	return &pgMetaStorage{db: db}
}
