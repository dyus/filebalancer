package storage

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type MetaStorage interface {
	Save(string, int64, []FilePart) error
	Get(string) (*fileMeta, error)
}

type pgMetaStorage struct {
	db *sqlx.DB
}

type FilePart struct {
	Path          string `db:"path"`
	ContentLength int64  `db:"content_length"`
}

type FilePartList []FilePart
type fileMeta struct {
	Id            int          `db:"id"`
	Name          string       `db:"name"`
	ContentLength int64        `db:"content_length"`
	FileParts     FilePartList `db:"parts"`
}

func (f FilePart) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FilePart) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), &f)
}

func (f FilePartList) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FilePartList) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), &f)
}

func (m *pgMetaStorage) Save(path string, contentLength int64, parts []FilePart) error {
	fileMeta := fileMeta{Name: path, ContentLength: contentLength, FileParts: parts}
	_, err := m.db.Exec(`INSERT INTO file_meta (name, content_length, parts)
		VALUES ($1, $2, $3)`,
		fileMeta.Name, fileMeta.ContentLength, fileMeta.FileParts)
	return err
}

func (m *pgMetaStorage) Get(name string) (*fileMeta, error) {
	file_meta := fileMeta{}
	err := m.db.Get(&file_meta, "SELECT * FROM file_meta WHERE name=$1", name)
	if err != nil {
		return nil, err
	}
	return &file_meta, nil
}

func NewMetaStorage(db *sqlx.DB) MetaStorage {
	return &pgMetaStorage{db: db}
}
