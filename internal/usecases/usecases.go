package usecases

import (
	"password-saver/internal/config"
	"password-saver/internal/core"
	"password-saver/internal/infrastructure/dto"
)

type UseCases struct {
	UserUseCase     *UserUseCase
	PasswordUseCase *PasswordUseCase
	SystemUseCase   *SystemUseCase
}

func InitUseCases(ur UserRepository, pr PasswordRepository, sr SystemRepository, cfg *config.EncryptKeys) *UseCases {
	return &UseCases{
		UserUseCase:     newUserUseCase(ur),
		PasswordUseCase: newPasswordUseCase(pr, cfg),
		SystemUseCase:   newSystemUseCase(sr),
	}
}

type PasswordRepository interface {
	Save(p *core.Password) error
	GetAll(userID int64) ([]core.Password, error)
	GetByID(passwordID, userID int64) (*core.Password, error)
	Update(p *core.Password) error
	Delete(passwordID, userID int64) error
}

type UserRepository interface {
	Registration(u *core.User) (int64, error)
	LogIn(q *dto.AuthRequest) (*core.User, error)
	Update(u *core.User) error
	Delete(userID int64) error
	GetByID(userID int64) (*core.User, error)
}

type SystemRepository interface {
	PingDB() error
}
