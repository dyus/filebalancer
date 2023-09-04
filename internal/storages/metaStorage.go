package storages

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type FileMeta struct {
	ContentLength int64  `json:"content_length"`
	Number        int    `json:"number"`
	StorageUrl    string `json:"storage_url"`
}

type FileMetaArr []FileMeta

type File struct {
	Id            int         `db:"id"`
	Name          string      `db:"name"`
	ContentLength int64       `db:"content_length"`
	FileMeta      FileMetaArr `db:"meta"`
}

type FileStorage struct {
	Db *sqlx.DB
}

func (m *FileMetaArr) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *FileMetaArr) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &m)
}

func (s *FileStorage) Add(file *File) error {
	meta, error := json.Marshal(file.FileMeta)
	if error != nil {
		return error
	}
	result, err := s.Db.Exec(
		`INSERT INTO file (name, content_length, meta) VALUES ($1, $2, $3)`,
		file.Name,
		file.ContentLength,
		meta,
	)
	if err != nil {
		return err
	}
	log.Info().Msg(fmt.Sprintf("call add %v+", result))
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStorage) Get(name string) (*File, error) {
	d := File{}
	error := s.Db.Get(&d, `select * from file where name=$1`, name)
	if error != nil {
		return nil, error
	}
	return &d, nil
}
