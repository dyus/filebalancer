package storage

import (
	"github.com/jmoiron/sqlx"
)

type MetaStorage interface {
	Save(string)
}

type pgMetaStorage struct {
	db *sqlx.DB
}

type filePart struct {
	Path          string `db:"path"`
	ContentLength int64  `db:"content_length"`
}

type fileMeta struct {
	Name string `db:"name"`
	// TODO read and impl converter https://jmoiron.github.io/sqlx/#advancedScanning
	FileParts []filePart
}

func (m *pgMetaStorage) Save(path string) {
	return
}

func NewMetaStorage(db *sqlx.DB) MetaStorage {
	return &pgMetaStorage{db: db}
}
