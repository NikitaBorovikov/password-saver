package model

import "password-saver/pkg/dto"

type (
	Password struct {
		PasswordID  int64
		UserID      int64
		Service     string
		EncPassword string
	}

	User struct {
		UserID       int64  `db:"user_id"`
		Email        string `db:"email"`
		HashPassword string `db:"hash_password"`
		RegDate      string `db:"reg_date"`
	}

	PasswordRepository interface {
		Save(p *Password) error
		GetAll(userID int64) ([]dto.PasswordResponse, error)
		GetByID(passwordID string) (*dto.PasswordResponse, error)
		Update(p *Password) error
		Delete(passwordID string) error
		//Generate New
	}

	UserRepository interface {
		Registration(u *User) (int64, error)
		LogIn(q *dto.LogInRequest) (*User, error)
		Update(u *User) error
		Delete(userID int64) error
		GetUserByID(userID int64) (*User, error)
	}
)
