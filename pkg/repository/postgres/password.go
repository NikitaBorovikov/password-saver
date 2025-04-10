package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	apperrors "password-saver/pkg/errors"
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
	_, err := r.db.NamedQuery(queryInserNewPassword, p)
	if err != nil {
		return fmt.Errorf("failed to save new password: %v", err)
	}

	return nil
}

func (r *PasswordRepository) GetAll(userID int64) ([]model.Password, error) {
	var passwords []model.Password

	if err := r.db.Select(&passwords, querySelectUserPasswords, userID); err != nil {
		return nil, fmt.Errorf("failed to select user's passwords: %v", err)
	}

	return passwords, nil
}

func (r *PasswordRepository) GetByID(passwordID int64) (*model.Password, error) {
	var password model.Password
	if err := r.db.Get(&password, querySelectPasswordByID, passwordID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrPasswordNotExists
		}

		return nil, fmt.Errorf("failed to get password by ID: %v", err)
	}

	return &password, nil
}

func (r *PasswordRepository) Update(p *model.Password) error {
	_, err := r.db.NamedExec(queryUpdatePassword, p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrPasswordNotExists
		}

		return fmt.Errorf("failed to update password: %v", err)
	}
	return nil
}

func (r *PasswordRepository) Delete(passwordID int64) error {
	_, err := r.db.Exec(queryDelPassword, passwordID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrPasswordNotExists
		}
		return fmt.Errorf("failed to delete password: %v", err)
	}
	return nil
}
