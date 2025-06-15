package db

import (
	"fmt"
	"password-saver/pkg/infrastructure/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	driverName = "postgres"
)

func ConnAndPing(cfg config.Postgres) (*sqlx.DB, error) {
	connectionStr := makeConnectionString(cfg)

	db, err := sqlx.Open(driverName, connectionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %v", err)
	}

	return db, nil
}

func makeConnectionString(cfg config.Postgres) string {
	connectionStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%d sslmode=disable",
		cfg.User, cfg.Name, cfg.Password, cfg.Host, cfg.Port,
	)
	return connectionStr
}
