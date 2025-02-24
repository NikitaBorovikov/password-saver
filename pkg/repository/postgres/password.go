package postgres

import (
	"password-saver/pkg/dto"
	"password-saver/pkg/model"

	"github.com/jmoiron/sqlx"
)

type PasswordRepository struct {
	db *sqlx.DB
}

func NewPasswordRepository(db *sqlx.DB) model.PasswordRepository {
	return &PasswordRepository{
		db: db,
	}
}

func (r *PasswordRepository) Save(p *model.Password) error {
	return nil
}

func (r *PasswordRepository) GetAll(userID int64) ([]dto.PasswordResponse, error) {
	return nil, nil
}

func (r *PasswordRepository) GetByID(passwordID string) (*dto.PasswordResponse, error) {
	return nil, nil
}

func (r *PasswordRepository) Update(p *model.Password) error {
	return nil
}

func (r *PasswordRepository) Delete(passwordID string) error {
	return nil
}
