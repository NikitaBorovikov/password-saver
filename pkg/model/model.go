package model

import "password-saver/pkg/dto"

type (
	Password struct {
		PasswordID  int64  `json:"password_id"`
		UserID      int64  `json:"user_id"`
		Service     string `json:"service"`
		EncPassword string `json:"-"`
	}

	User struct {
		UserID       int64  `json:"user_id"`
		Email        string `json:"email" validate:"required,email"`
		Password     string `json:"-" validate:"min=7,max=255"`
		HashPassword string `json:"-"`
		Salt         string `json:"-"`
		RegDate      string `json:"reg_date"`
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
		Registration(u *User) error
		LogIn(q *dto.LogInRequest) (*User, error)
		Update(u *User) error
		Delete(userID int64) error
	}
)
