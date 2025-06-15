package postgres

import (
	"password-saver/pkg/core"

	"github.com/jmoiron/sqlx"
)

type SystemRepository struct {
	db *sqlx.DB
}

func NewSystemRepository(db *sqlx.DB) core.SystemRepository {
	return &SystemRepository{
		db: db,
	}
}

func (r *SystemRepository) PingDB() error {
	if err := r.db.DB.Ping(); err != nil {
		return ErrPingFailed
	}
	return nil
}
