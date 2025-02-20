package api

import "password-saver/pkg/usecases"

type Handlers struct {
	userHandler     *UserHandler
	passwordHandler *PasswordHandler
}

func InitHandlers(uc *usecases.UseCases) *Handlers {
	return &Handlers{
		userHandler:     newUserHandler(uc.UserUseCase),
		passwordHandler: newPasswordHandler(uc.PasswordUseCase),
	}
}
