package model

import "password-saver/pkg/dto"

type (
	Password struct {
		PasswordID  int64   `db:"password_id"`
		UserID      int64   `db:"user_id"`
		EncService  string  `db:"enc_service"`
		EncPassword string  `db:"enc_password"`
		EncLogin    *string `db:"enc_login"`
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
		GetByID(passwordID int64) (*Password, error)
		Update(p *Password) error
		Delete(passwordID int64) error
	}

	UserRepository interface {
		Registration(u *User) (int64, error)
		LogIn(q *dto.AuthRequest) (*User, error)
		Update(u *User) error
		Delete(userID int64) error
		GetByID(userID int64) (*User, error)
	}
)
