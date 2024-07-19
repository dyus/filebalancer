package app

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func newDb(conf *DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", conf.Username, conf.Password, conf.Name, conf.Host, conf.Port))
	if err != nil {
		return nil, err
	}
	sql, err := os.ReadFile("db/create_db.sql")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		return nil, err
	}
	return db, nil
}
