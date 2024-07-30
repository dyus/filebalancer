package app

import (
	"github.com/dyus/filebalancer/internal/storage"
	"github.com/jmoiron/sqlx"
)

func NewMetaStorage(db *sqlx.DB) storage.MetaStorage {
	return storage.NewMetaStorage(db)
}
