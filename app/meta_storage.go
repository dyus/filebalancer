package app

import (
	"github.com/dyus/filebalancer/internal/storage"
	"github.com/jmoiron/sqlx"
)

func newMetaStorage(db *sqlx.DB) storage.MetaStorage {
	return storage.NewMetaStorage(db)
}
