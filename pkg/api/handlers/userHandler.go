package api

import (
	"password-saver/pkg/usecases"
)

type UserHandler struct {
	UserHandler *usecases.UserUseCase
}

func newUserHandler(uc *usecases.UserUseCase) *UserHandler {
	return &UserHandler{
		UserHandler: uc,
	}
}
