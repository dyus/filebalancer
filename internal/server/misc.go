package server

import (
	"filebalancer/internal/storages"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func SampleDbQuering(db *sqlx.DB) {

	metaStorage := storages.FileStorage{Db: db}
	meta := storages.FileMetaArr{
		{
			ContentLength: 100,
			Number:        1,
			StorageUrl:    "test/123",
		},
	}
	file := &storages.File{
		Name:          "test",
		ContentLength: 10,
		FileMeta:      meta,
	}
	if err := metaStorage.Add(file); err != nil {
		panic(err)
	}
	res, err := metaStorage.Get(file.Name)
	if err != nil {
		panic(err)
	}
	log.Info().Msg(fmt.Sprintf("%v+", res))

}
