package usecases

import (
	"password-saver/pkg/config"
	"password-saver/pkg/repository"
)

type UseCases struct {
	UserUseCase     *UserUseCase
	PasswordUseCase *PasswordUseCase
	SystemUseCase   *SystemUseCase
}

func InitUseCases(r *repository.Repository, cfg *config.EncryptKeys) *UseCases {
	return &UseCases{
		UserUseCase:     newUserUseCase(r.UserRepository),
		PasswordUseCase: newPasswordUseCase(r.PasswordRepository, cfg),
		SystemUseCase:   newSystemUseCase(r.SystemRepository),
	}
}
