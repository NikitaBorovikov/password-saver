package postgres

import (
	"password-saver/pkg/core"
	"password-saver/pkg/dto"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) core.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Registration(u *core.User) (int64, error) {
	var userID int64

	rows, err := r.db.NamedQuery(queryRegistration, u)
	if err != nil {
		return 0, handleSQLErrors(err)
	}

	if rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			return 0, ErrScanFailed
		}
	}
	return userID, nil
}

func (r *UserRepository) LogIn(q *dto.AuthRequest) (*core.User, error) {
	var user core.User

	if err := r.db.Get(&user, queryLogIn, q.Email); err != nil {
		return nil, handleSQLErrors(err)
	}
	return &user, nil
}

func (r *UserRepository) Update(u *core.User) error {
	_, err := r.db.NamedExec(queryUpdateUser, u)
	if err != nil {
		return handleSQLErrors(err)
	}
	return nil
}

func (r *UserRepository) Delete(userID int64) error {
	_, err := r.db.Exec(queryDelUser, userID)
	if err != nil {
		return handleSQLErrors(err)
	}

	return nil
}

func (r *UserRepository) GetByID(userID int64) (*core.User, error) {
	var user core.User

	if err := r.db.Get(&user, querySelectUserByID, userID); err != nil {
		return nil, handleSQLErrors(err)
	}
	return &user, nil
}
