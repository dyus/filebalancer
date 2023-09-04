package server

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func NewDb() (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", "user=postgres password=1 dbname=filebalancer sslmode=disable")
	if err != nil {
		log.Error().Err(err).Msg("Can't connect sqlite3")
		return nil, err
	}
	return conn, nil
}

func SetupDb(db *sqlx.DB) {
	var schema = `--sql
DROP TABLE IF EXISTS file;
CREATE TABLE file
(
	id 				SERIAL PRIMARY KEY,
	name            VARCHAR(1024),
    content_length  BIGINT,
	meta			JSON
);

-- CREATE TABLE meta
-- (
-- 	id				BIGINT AUTO_INCREMENT PRIMARY KEY,
-- 	content_length  BIGINT,
-- 	number			integer,
-- 	file_id 		BIGINT,
-- 	storage_url		TEXT,
-- 	FOREIGN KEY (file_id) REFERENCES file(id)		
-- )
`
	db.MustExec(schema)
}
