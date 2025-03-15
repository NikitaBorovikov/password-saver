package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"password-saver/pkg/dto"
	apperrors "password-saver/pkg/errors"
	"password-saver/pkg/model"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	uniqueViolationErrCode = "23505"
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
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == uniqueViolationErrCode {
				return 0, apperrors.ErrDuplicateUser
			}
		}
		return 0, fmt.Errorf("registration error db: %v", err)
	}

	if rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			return 0, fmt.Errorf("failed to scan user ID: %v", err)
		}
	}
	return userID, nil
}

func (r *UserRepository) LogIn(q *dto.AuthRequest) (*model.User, error) {
	var user model.User

	if err := r.db.Get(&user, queryLogIn, q.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to login user: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(u *model.User) error {
	_, err := r.db.NamedExec(queryUpdateUser, u)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrUserNotFound
		}
		return fmt.Errorf("failed to update user: %v", err)

	}
	return nil
}

func (r *UserRepository) Delete(userID int64) error {
	_, err := r.db.Exec(queryDelUser, userID)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrUserNotFound
		}
		return fmt.Errorf("failed to delete user: %v", err)

	}
	//TODO: delete users' passwords from passwords table
	return nil
}

func (r *UserRepository) GetByID(userID int64) (*model.User, error) {
	var user model.User
	if err := r.db.Get(&user, queryGetUserByID, userID); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by ID: %v", err)

	}

	return &user, nil
}
