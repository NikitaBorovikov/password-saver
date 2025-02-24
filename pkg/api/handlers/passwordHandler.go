package handlers

import "password-saver/pkg/usecases"

type PasswordHandler struct {
	PasswordUseCase *usecases.PasswordUseCase
}

func newPasswordHandler(uc *usecases.PasswordUseCase) *PasswordHandler {
	return &PasswordHandler{
		PasswordUseCase: uc,
	}
}
