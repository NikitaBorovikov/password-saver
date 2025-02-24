package handlers

import (
	"net/http"
	"password-saver/pkg/dto"
	"password-saver/pkg/model"
	"password-saver/pkg/usecases"

	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
)

type UserHandler struct {
	UserUseCase *usecases.UserUseCase
}

func newUserHandler(uc *usecases.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: uc,
	}
}

func (h *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {
	req := dto.RegRequest{}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	user := model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.UserUseCase.Registration(&user); err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	sendOKResponse(w, r, req.Email)
}

func (h *UserHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	req := dto.LogInRequest{}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := h.UserUseCase.LogIn(&req)
	if err != nil {
		sendErrorRespose(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := setUserSession(w, r, user); err != nil {
		sendErrorRespose(w, r, http.StatusInternalServerError, err)
		return
	}

	sendOKResponse(w, r, user.UserID)
}

func setUserSession(w http.ResponseWriter, r *http.Request, user *model.User) error {
	session, err := sessionStore.Get(r, sessionKey)
	if err != nil {
		return err
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = user.UserID

	session.Options = &sessions.Options{
		MaxAge:   3600 * 12,
		HttpOnly: true,
	}

	err = session.Save(r, w)
	return err
}
