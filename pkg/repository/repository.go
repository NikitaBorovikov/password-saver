package repository

import (
	"password-saver/pkg/model"
	"password-saver/pkg/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	UserRepository     model.UserRepository
	PasswordRepository model.PasswordRepository
}

func InitRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:     postgres.NewUserRepository(db),
		PasswordRepository: postgres.NewPasswordRepository(db),
	}
}
