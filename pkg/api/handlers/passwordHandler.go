package handlers

import (
	"fmt"
	"net/http"
	"password-saver/pkg/dto"
	"password-saver/pkg/usecases"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type PasswordHandler struct {
	PasswordUseCase *usecases.PasswordUseCase
}

func newPasswordHandler(uc *usecases.PasswordUseCase) *PasswordHandler {
	return &PasswordHandler{
		PasswordUseCase: uc,
	}
}

func (h *PasswordHandler) Save(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		err := fmt.Errorf("no userID in context")
		sendErrorRespose(w, r, http.StatusUnauthorized, err)
		return
	}

	req, err := decodeSavePasswordRequest(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.PasswordUseCase.Save(req, userID); err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}
}

func (h *PasswordHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		err := fmt.Errorf("no userID in context")
		sendErrorRespose(w, r, http.StatusUnauthorized, err)
		return
	}

	userPasswords, err := h.PasswordUseCase.GetAll(userID)
	if err != nil {
		sendErrorRespose(w, r, http.StatusUnauthorized, err)
		return
	}

	sendOKResponse(w, r, userPasswords)
}

func (h *PasswordHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	passwordResponse, err := h.PasswordUseCase.GetByID(passwordID)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	sendOKResponse(w, r, passwordResponse)
}

func (h *PasswordHandler) Update(w http.ResponseWriter, r *http.Request) {
	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		err := fmt.Errorf("no userID in context")
		sendErrorRespose(w, r, http.StatusUnauthorized, err)
		return
	}

	req, err := decodeUpdatePasswordRequest(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.PasswordUseCase.Update(req, passwordID, userID); err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	sendOKResponse(w, r, "password is updated")
}

func (h *PasswordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.PasswordUseCase.Delete(passwordID); err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	sendOKResponse(w, r, nil)
}

func decodeSavePasswordRequest(r *http.Request) (*dto.SavePasswordRequest, error) {
	var req dto.SavePasswordRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %v", err)
	}

	return &req, nil
}

func decodeUpdatePasswordRequest(r *http.Request) (*dto.UpdatePasswordRequest, error) {
	var req dto.UpdatePasswordRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %v", err)
	}

	return &req, nil
}

func getPasswordIDFromURL(r *http.Request) (int64, error) {
	passwordID := chi.URLParam(r, "passwordID")

	passwordIDInt, err := strconv.Atoi(passwordID)
	if err != nil {
		return 0, fmt.Errorf("failed to convert str to int")
	}

	return int64(passwordIDInt), nil
}
