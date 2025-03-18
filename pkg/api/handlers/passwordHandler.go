package handlers

import (
	"net/http"
	"password-saver/pkg/dto"
	apperrors "password-saver/pkg/errors"
	"password-saver/pkg/usecases"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
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
		logrus.Error("failed to get userID from context")
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	req, err := decodePasswordRequest(r)
	if err != nil {
		logrus.Errorf("failed to decode request: %v", err)
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrDecodeRequest)
		return
	}

	if err := h.PasswordUseCase.Save(req, userID); err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	sendOKResponse(w, r, http.StatusCreated, "password is saved")

	logrus.Info("password was saved sucessfully")
}

func (h *PasswordHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error("failed to get userID from context")
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	userPasswords, err := h.PasswordUseCase.GetAll(userID)
	if err != nil {
		sendErrorRespose(w, r, http.StatusInternalServerError, err)
		return
	}

	sendOKResponse(w, r, http.StatusOK, userPasswords)

	logrus.Info("passwords was given sucessfully")
}

func (h *PasswordHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		logrus.Error("failed to get passwordID from url")
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrInvalidURLParam)
		return
	}

	passwordResponse, err := h.PasswordUseCase.GetByID(passwordID)
	if err != nil {
		sendErrorRespose(w, r, http.StatusInternalServerError, err)
		return
	}

	sendOKResponse(w, r, http.StatusOK, passwordResponse)

	logrus.Info("password was given sucessfully")
}

func (h *PasswordHandler) Update(w http.ResponseWriter, r *http.Request) {
	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error("failed to get userID from context")
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	req, err := decodePasswordRequest(r)
	if err != nil {
		logrus.Errorf("failed to decode request: %v", err)
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrDecodeRequest)
		return
	}

	if err := h.PasswordUseCase.Update(req, passwordID, userID); err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	sendOKResponse(w, r, http.StatusOK, "password is updated")

	logrus.Info("passwords was updated sucessfully")
}

func (h *PasswordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		logrus.Error("failed to get passwordID from url")
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrInvalidURLParam)
		return
	}

	if err := h.PasswordUseCase.Delete(passwordID); err != nil {
		sendErrorRespose(w, r, http.StatusInternalServerError, err)
		return
	}

	sendOKResponse(w, r, http.StatusNoContent, "password is deleted")

	logrus.Info("passwords was deleted sucessfully")
}

func (h *PasswordHandler) Generate(w http.ResponseWriter, r *http.Request) {
	ps, err := getPasswordSettingsFromURL(r)
	if err != nil {
		logrus.Errorf("failed to get password setting for geneating from URL: %v", err)
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrInvalidURLParam)
		return
	}

	password, err := h.PasswordUseCase.Generate(ps.length, ps.useSpecialSymbols)
	if err != nil {
		sendErrorRespose(w, r, http.StatusInternalServerError, err)
	}

	sendOKResponse(w, r, http.StatusOK, password)

	logrus.Info("new password was generated successfully")

}

type passwordSettings struct {
	length            int
	useSpecialSymbols bool
}

func getPasswordSettingsFromURL(r *http.Request) (*passwordSettings, error) {
	lenStr := r.URL.Query().Get("len")
	len, err := strconv.Atoi(lenStr)
	if err != nil {
		return nil, err
	}

	useSpecialStr := r.URL.Query().Get("special")
	useSpecial := useSpecialStr == "true"

	ps := &passwordSettings{
		length:            len,
		useSpecialSymbols: useSpecial,
	}

	return ps, nil
}

func decodePasswordRequest(r *http.Request) (*dto.PasswordRequest, error) {
	var req dto.PasswordRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func getPasswordIDFromURL(r *http.Request) (int64, error) {
	passwordID := chi.URLParam(r, "passwordID")

	passwordIDInt, err := strconv.Atoi(passwordID)
	if err != nil {
		return 0, err
	}

	return int64(passwordIDInt), nil
}
