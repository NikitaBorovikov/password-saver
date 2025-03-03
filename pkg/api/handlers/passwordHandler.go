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
	req, err := decodeSavePasswordRequest(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.PasswordUseCase.Save(req); err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}
}

func decodeSavePasswordRequest(r *http.Request) (*dto.SavePasswordRequest, error) {
	var req dto.SavePasswordRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %v", err)
	}

	return &req, nil
}
