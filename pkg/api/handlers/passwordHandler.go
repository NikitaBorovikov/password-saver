package handlers

import "password-saver/pkg/usecases"

type PasswordHandler struct {
	PasswordHandler *usecases.PasswordUseCase
}

func newPasswordHandler(uc *usecases.PasswordUseCase) *PasswordHandler {
	return &PasswordHandler{
		PasswordHandler: uc,
	}
}
