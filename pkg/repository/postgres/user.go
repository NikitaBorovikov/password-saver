package postgres

import (
	"fmt"
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
		return nil, fmt.Errorf("failed to login user: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(u *model.User) error {
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

func (r *UserRepository) GetUserByID(userID int64) (*model.User, error) {
	return nil, nil
}
