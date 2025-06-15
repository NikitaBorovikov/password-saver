package repository

import (
	"password-saver/pkg/core"
	"password-saver/pkg/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	UserRepository     core.UserRepository
	PasswordRepository core.PasswordRepository
	SystemRepository   core.SystemRepository
}

func InitRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:     postgres.NewUserRepository(db),
		PasswordRepository: postgres.NewPasswordRepository(db),
		SystemRepository:   postgres.NewSystemRepository(db),
	}
}
