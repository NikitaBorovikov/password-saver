package postgres

import (
	"password-saver/pkg/dto"
	"password-saver/pkg/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Registration(u *model.User) error {
	return nil
}

func (r *UserRepository) LogIn(q *dto.LogInRequest) (*model.User, error) {
	return nil, nil
}

func (r *UserRepository) Update(u *model.User) error {
	return nil
}

func (r *UserRepository) Delete(userID int64) error {
	return nil
}
