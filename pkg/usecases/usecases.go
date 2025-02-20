package usecases

import "password-saver/pkg/repository"

type UseCases struct {
	UserUseCase     *UserUseCase
	PasswordUseCase *PasswordUseCase
}

func InitUseCases(r *repository.Repository) *UseCases {
	return &UseCases{
		UserUseCase:     NewUserUseCase(r.UserRepository),
		PasswordUseCase: NewPasswordUseCase(r.PasswordRepository),
	}
}
