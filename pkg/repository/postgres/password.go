package postgres

import (
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
		return handleSQLErrors(err)
	}
	return nil
}

func (r *PasswordRepository) GetAll(userID int64) ([]model.Password, error) {
	var passwords []model.Password

	if err := r.db.Select(&passwords, querySelectUserPasswords, userID); err != nil {
		return nil, handleSQLErrors(err)
	}
	return passwords, nil
}

func (r *PasswordRepository) GetByID(passwordID, userID int64) (*model.Password, error) {
	var password model.Password

	if err := r.db.Get(&password, querySelectPasswordByID, passwordID, userID); err != nil {
		return nil, handleSQLErrors(err)
	}
	return &password, nil
}

func (r *PasswordRepository) Update(p *model.Password) error {
	_, err := r.db.NamedExec(queryUpdatePassword, p)
	if err != nil {
		return handleSQLErrors(err)
	}
	return nil
}

func (r *PasswordRepository) Delete(passwordID, userID int64) error {
	_, err := r.db.Exec(queryDelPassword, passwordID, userID)
	if err != nil {
		return handleSQLErrors(err)
	}
	return nil
}
