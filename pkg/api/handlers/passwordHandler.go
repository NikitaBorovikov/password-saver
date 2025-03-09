package handlers

import (
	"fmt"
	"net/http"
	"password-saver/pkg/dto"
	"password-saver/pkg/usecases"

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

func decodeSavePasswordRequest(r *http.Request) (*dto.SavePasswordRequest, error) {
	var req dto.SavePasswordRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %v", err)
	}

	return &req, nil
}
