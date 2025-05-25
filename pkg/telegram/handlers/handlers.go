package handlers

import "password-saver/pkg/usecases"

type Handlers struct {
	PasswordHandler PasswordHandler
}

func InitHandlers(uc *usecases.UseCases) *Handlers {
	return &Handlers{
		PasswordHandler: *newPasswordHandler(uc.PasswordUseCase),
	}
}
