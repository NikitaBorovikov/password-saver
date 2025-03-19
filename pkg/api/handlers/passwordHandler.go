package handlers

import (
	"net/http"
	"password-saver/pkg/dto"
	apperrors "password-saver/pkg/errors"
	"password-saver/pkg/logs"
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
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	req, err := decodePasswordRequest(r)
	if err != nil {
		logrus.Errorf(logs.FailedToDecodeRequest, err)
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrDecodeRequest)
		return
	}

	if err := h.PasswordUseCase.Save(req, userID); err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	sendOKResponse(w, r, http.StatusCreated, nil)

	logrus.Info(logs.PasswordSavedSuccessfully)
}

func (h *PasswordHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	userPasswords, err := h.PasswordUseCase.GetAll(userID)
	if err != nil {
		sendErrorRespose(w, r, http.StatusInternalServerError, err)
		return
	}

	sendOKResponse(w, r, http.StatusOK, userPasswords)

	logrus.Info(logs.PasswordsGivenSuccessfully)
}

func (h *PasswordHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		logrus.Error(logs.FailedToGetPasswordIDFromURL)
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrInvalidURLParam)
		return
	}

	passwordResponse, err := h.PasswordUseCase.GetByID(passwordID)
	if err != nil {
		sendErrorRespose(w, r, http.StatusInternalServerError, err)
		return
	}

	sendOKResponse(w, r, http.StatusOK, passwordResponse)

	logrus.Info(logs.PasswordsGivenSuccessfully)
}

func (h *PasswordHandler) Update(w http.ResponseWriter, r *http.Request) {
	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	req, err := decodePasswordRequest(r)
	if err != nil {
		logrus.Errorf(logs.FailedToDecodeRequest, err)
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrDecodeRequest)
		return
	}

	if err := h.PasswordUseCase.Update(req, passwordID, userID); err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	sendOKResponse(w, r, http.StatusOK, nil)

	logrus.Info(logs.PasswordUpdatedSuccessfully)
}

func (h *PasswordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		logrus.Error(logs.FailedToGetPasswordIDFromURL)
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrInvalidURLParam)
		return
	}

	if err := h.PasswordUseCase.Delete(passwordID); err != nil {
		sendErrorRespose(w, r, http.StatusInternalServerError, err)
		return
	}

	sendOKResponse(w, r, http.StatusNoContent, nil)

	logrus.Info(logs.PasswordDeletedSuccessfully)
}

func (h *PasswordHandler) Generate(w http.ResponseWriter, r *http.Request) {
	ps, err := getPasswordSettingsFromURL(r)
	if err != nil {
		logrus.Errorf(logs.FailedToGetPasswordSettings, err)
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrInvalidURLParam)
		return
	}

	password, err := h.PasswordUseCase.Generate(ps)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	sendOKResponse(w, r, http.StatusOK, password)

	logrus.Info(logs.PasswordGeneratedSuccessfully)
}

func getPasswordSettingsFromURL(r *http.Request) (*dto.GeneratePasswordRequest, error) {
	lenStr := r.URL.Query().Get("len")
	len, err := strconv.Atoi(lenStr)
	if err != nil {
		return nil, err
	}

	useSpecialStr := r.URL.Query().Get("special")
	useSpecial := useSpecialStr == "true"

	ps := &dto.GeneratePasswordRequest{
		Length:            len,
		UseSpecialSymbols: useSpecial,
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
