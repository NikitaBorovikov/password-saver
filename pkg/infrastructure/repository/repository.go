package repository

import (
	"password-saver/pkg/infrastructure/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	UserRepository     *postgres.UserRepository
	PasswordRepository *postgres.PasswordRepository
	SystemRepository   *postgres.SystemRepository
}

func InitRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:     postgres.NewUserRepository(db),
		PasswordRepository: postgres.NewPasswordRepository(db),
		SystemRepository:   postgres.NewSystemRepository(db),
	}
}
