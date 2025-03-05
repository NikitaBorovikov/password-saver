package handlers

import (
	"context"
	"net/http"
	"password-saver/pkg/api/session"
	"password-saver/pkg/dto"
	"password-saver/pkg/usecases"

	"github.com/go-chi/render"
)

type Handlers struct {
	UserHandler     *UserHandler
	PasswordHandler *PasswordHandler
}

func InitHandlers(uc *usecases.UseCases, session *session.SessionManager) *Handlers {
	return &Handlers{
		UserHandler:     newUserHandler(uc.UserUseCase, session),
		PasswordHandler: newPasswordHandler(uc.PasswordUseCase),
	}
}

func getUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

func sendErrorRespose(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.WriteHeader(statusCode)
	render.JSON(w, r, dto.NewErrorResponse(err))
}

func sendOKResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, data)
}
