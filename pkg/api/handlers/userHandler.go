package handlers

import (
	"fmt"
	"net/http"
	"password-saver/pkg/dto"
	"password-saver/pkg/model"
	"password-saver/pkg/usecases"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
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

	req, err := decodeRegRequest(r)
	if err != nil {
		logrus.Error(err)
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	user := model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.UserUseCase.Registration(&user); err != nil {
		logrus.Error(err)
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	sendOKResponse(w, r, req.Email)
}

func (h *UserHandler) LogIn(w http.ResponseWriter, r *http.Request) {

	req, err := decodeLogInRequest(r)
	if err != nil {
		logrus.Error(err)
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := h.UserUseCase.LogIn(req)
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

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		err := fmt.Errorf("no userID in context")
		sendErrorRespose(w, r, http.StatusUnauthorized, err)
		return
	}

	req, err := decodeUpdateRequest(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	req.UserID = userID

	if err := h.UserUseCase.Update(req); err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	sendOKResponse(w, r, nil)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

func decodeRegRequest(r *http.Request) (*dto.RegRequest, error) {
	var req dto.RegRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %v", err)
	}

	return &req, nil
}

func decodeLogInRequest(r *http.Request) (*dto.LogInRequest, error) {
	var req dto.LogInRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %v", err)
	}

	return &req, nil
}

func decodeUpdateRequest(r *http.Request) (*dto.UpdateUserRequest, error) {
	var req dto.UpdateUserRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %v", err)
	}

	return &req, nil
}
