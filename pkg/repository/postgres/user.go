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

func (r *UserRepository) Registration(u *model.User) (int64, error) {
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

func (r *UserRepository) LogIn(q *dto.AuthRequest) (*model.User, error) {
	var user model.User

	if err := r.db.Get(&user, queryLogIn, q.Email); err != nil {
		return nil, handleSQLErrors(err)
	}
	return &user, nil
}

func (r *UserRepository) Update(u *model.User) error {
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

func (r *UserRepository) GetByID(userID int64) (*model.User, error) {
	var user model.User

	if err := r.db.Get(&user, querySelectUserByID, userID); err != nil {
		return nil, handleSQLErrors(err)
	}
	return &user, nil
}
