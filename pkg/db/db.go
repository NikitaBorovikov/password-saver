package db

import (
	"password-saver/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnAndPing(cfg config.Postgres) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	return db, err

}
