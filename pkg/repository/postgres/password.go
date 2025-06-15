package postgres

import (
	"password-saver/pkg/core"

	"github.com/jmoiron/sqlx"
)

type PasswordRepository struct {
	db *sqlx.DB
}

func NewPasswordRepository(db *sqlx.DB) core.PasswordRepository {
	return &PasswordRepository{
		db: db,
	}
}

func (r *PasswordRepository) Save(p *core.Password) error {
	_, err := r.db.NamedQuery(queryInserNewPassword, p)
	if err != nil {
		return handleSQLErrors(err)
	}
	return nil
}

func (r *PasswordRepository) GetAll(userID int64) ([]core.Password, error) {
	var passwords []core.Password

	if err := r.db.Select(&passwords, querySelectUserPasswords, userID); err != nil {
		return nil, handleSQLErrors(err)
	}
	return passwords, nil
}

func (r *PasswordRepository) GetByID(passwordID, userID int64) (*core.Password, error) {
	var password core.Password

	if err := r.db.Get(&password, querySelectPasswordByID, passwordID, userID); err != nil {
		return nil, handleSQLErrors(err)
	}
	return &password, nil
}

func (r *PasswordRepository) Update(p *core.Password) error {
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
