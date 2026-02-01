package handlers

import (
	"context"
	"net/http"
	"password-saver/internal/infrastructure/api/session"
	"password-saver/internal/infrastructure/dto"
	"password-saver/internal/usecases"

	"github.com/go-chi/render"
)

type Handlers struct {
	UserHandler     *UserHandler
	PasswordHandler *PasswordHandler
	SystemHandler   *SystemHandler
}

type contextKey string

const (
	UserIDCtx contextKey = "userID"
)

func InitHandlers(uc *usecases.UseCases, session *session.SessionManager) *Handlers {
	return &Handlers{
		UserHandler:     newUserHandler(uc.UserUseCase, session),
		PasswordHandler: newPasswordHandler(uc.PasswordUseCase),
		SystemHandler:   newSystemHandler(uc.SystemUseCase),
	}
}

func getUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDCtx).(int64)
	return userID, ok
}

func sendErrorRespose(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.WriteHeader(statusCode)
	render.JSON(w, r, dto.NewErrorResponse(err))
}

func sendOKResponse(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	render.JSON(w, r, data)
}
