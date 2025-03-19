package db

import (
	"fmt"
	"password-saver/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	driverName = "postgres"
)

func ConnAndPing(cfg config.Postgres) (*sqlx.DB, error) {

	db, err := sqlx.Open(driverName, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %v", err)
	}

	return db, nil
}
