package model

import "password-saver/pkg/dto"

type (
	Password struct {
		PasswordID  int64  `db:"password_id"`
		UserID      int64  `db:"user_id"`
		EncService  string `db:"enc_service"`
		EncPassword string `db:"enc_password"`
	}

	User struct {
		UserID       int64  `db:"user_id"`
		Email        string `db:"email"`
		HashPassword string `db:"hash_password"`
		RegDate      string `db:"reg_date"`
	}

	PasswordRepository interface {
		Save(p *Password) error
		GetAll(userID int64) ([]Password, error)
		GetByID(passwordID string) (*Password, error)
		Update(p *Password) error
		Delete(passwordID int64) error
		//Generate New
	}

	UserRepository interface {
		Registration(u *User) (int64, error)
		LogIn(q *dto.LogInRequest) (*User, error)
		Update(u *User) error
		Delete(userID int64) error
		GetByID(userID int64) (*User, error)
	}
)
