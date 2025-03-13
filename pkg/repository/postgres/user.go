package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"password-saver/pkg/dto"
	"password-saver/pkg/model"

	"github.com/jmoiron/sqlx"
)

var (
	errUserIDNotExists = errors.New("such id doesn't exists")
	errEmailNotExists  = errors.New("there is no user with this email address")
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
		return 0, fmt.Errorf("registration error db: %v", err)
	}

	if rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			return 0, fmt.Errorf("failed to scan user ID: %v", err)
		}
	}
	return userID, nil
}

func (r *UserRepository) LogIn(q *dto.LogInRequest) (*model.User, error) {
	var user model.User

	if err := r.db.Get(&user, queryLogIn, q.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, errEmailNotExists
		}
		return nil, fmt.Errorf("failed to login user: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(u *model.User) error {
	_, err := r.db.NamedExec(queryUpdateUser, u)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

func (r *UserRepository) Delete(userID int64) error {
	_, err := r.db.Exec(queryDelUser, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	//TODO: delete users' passwords from passwords table
	return nil
}

func (r *UserRepository) GetByID(userID int64) (*model.User, error) {
	var user model.User
	if err := r.db.Get(&user, queryGetUserByID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errUserIDNotExists
		}

		return nil, fmt.Errorf("failed to get user by ID: %v", err)
	}

	return &user, nil
}
