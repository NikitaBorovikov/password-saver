package model

import "password-saver/pkg/dto"

type Password struct {
	PasswordID  int64
	UserID      int64
	Service     string
	EncPassword string
	//Password
}

type User struct {
	UserID       int64
	Email        string
	HashPassword string
	Salt         string
}

type PasswordRepository interface {
	Save(p *Password) error
	GetAll(userID int64) ([]Password, error)
	GetByID(passwordID string) (*Password, error)
	Update(p *Password) error
	Delete(passwordID string) error
	//Generate New
}

type UserRepository interface {
	Registration(u *User) error
	LogIn(q *dto.RegRequest) (*User, error)
	Update(u *User) error
	Delete(userID int64) error
}
