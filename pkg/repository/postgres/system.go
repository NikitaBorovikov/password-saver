package postgres

import (
	"fmt"
	"password-saver/pkg/model"

	"github.com/jmoiron/sqlx"
)

type SystemRepository struct {
	db *sqlx.DB
}

func NewSystemRepository(db *sqlx.DB) model.SystemRepository {
	return &SystemRepository{
		db: db,
	}
}

func (r *SystemRepository) PingDB() error {
	if err := r.db.DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %v", err)
	}
	return nil
}
