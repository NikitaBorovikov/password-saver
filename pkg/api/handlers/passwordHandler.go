package handlers

import (
	"net/http"
	"password-saver/pkg/dto"
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

// @Summary Save password
// @Description Save user's password (an active session is required).
// @Tags Passwords
// @Accept json
// @Produce json
// @Param input body dto.PasswordRequest true "Password data"
// @Success 201
// @Failure 400,422,500 {object} dto.ErrorResponse
// @Security SessionCookie
// @Router /passwords/ [post]
func (h *PasswordHandler) Save(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusUnauthorized, ErrNotAuthenticated)
		return
	}

	req, err := decodePasswordRequest(r)
	if err != nil {
		logrus.Errorf(logs.FailedToDecodeRequest, err)
		sendErrorRespose(w, r, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	if err := h.PasswordUseCase.Save(req, userID); err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
		return
	}

	sendOKResponse(w, r, http.StatusCreated, nil)

	logrus.Info(logs.PasswordSavedSuccessfully)
}

// @Summary Get passwords
// @Description Get all user's passwords by user ID form context (an active session is required).
// @Tags Passwords
// @Produce json
// @Success 200 {object} []dto.PasswordResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security SessionCookie
// @Router /passwords/ [get]
func (h *PasswordHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusUnauthorized, ErrNotAuthenticated)
		return
	}

	userPasswords, err := h.PasswordUseCase.GetAll(userID)
	if err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
		return
	}

	sendOKResponse(w, r, http.StatusOK, userPasswords)

	logrus.Info(logs.PasswordsGivenSuccessfully)
}

// @Summary Get password by ID
// @Description Get user's passwords by user ID form context and passwordID from URL (an active session is required).
// @Tags Passwords
// @Produce json
// @Param passwordID path string true "password ID"
// @Success 200 {object} dto.PasswordResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security SessionCookie
// @Router /passwords/{passwordID} [get]
func (h *PasswordHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		logrus.Error(logs.FailedToGetPasswordIDFromURL)
		sendErrorRespose(w, r, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	passwordResponse, err := h.PasswordUseCase.GetByID(passwordID, userID)
	if err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
		return
	}

	sendOKResponse(w, r, http.StatusOK, passwordResponse)

	logrus.Info(logs.PasswordsGivenSuccessfully)
}

// @Summary Update password's data
// @Description Update passwords's data by user ID from context and password ID from URL (an active session is required).
// @Tags Passwords
// @Accept json
// @Produce json
// @Param passwordID path string true "password ID"
// @Param input body dto.PasswordRequest true "New password data"
// @Success 200
// @Failure 400,422,500 {object} dto.ErrorResponse
// @Security SessionCookie
// @Router /passwords/ [put]
func (h *PasswordHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusUnauthorized, ErrNotAuthenticated)
		return
	}

	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	req, err := decodePasswordRequest(r)
	if err != nil {
		logrus.Errorf(logs.FailedToDecodeRequest, err)
		sendErrorRespose(w, r, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	if err := h.PasswordUseCase.Update(req, passwordID, userID); err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
		return
	}

	sendOKResponse(w, r, http.StatusOK, nil)

	logrus.Info(logs.PasswordUpdatedSuccessfully)
}

// @Summary Delete password
// @Description delete user's password by user id from context and password ID from URL (an active session is required).
// @Tags Passwords
// @Produce json
// @Param passwordID path string true "password ID"
// @Success 204
// @Failure 400,500 {object} dto.ErrorResponse
// @Security SessionCookie
// @Router /passwords/ [delete]
func (h *PasswordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusUnauthorized, ErrNotAuthenticated)
		return
	}

	passwordID, err := getPasswordIDFromURL(r)
	if err != nil {
		logrus.Error(logs.FailedToGetPasswordIDFromURL)
		sendErrorRespose(w, r, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	if err := h.PasswordUseCase.Delete(passwordID, userID); err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
		return
	}

	sendOKResponse(w, r, http.StatusNoContent, nil)

	logrus.Info(logs.PasswordDeletedSuccessfully)
}

// @Summary Generate password
// @Description Generate a new password based on the specified parameters (length, using special symbols)
// @Tags Passwords
// @Produce json
// @Param len query string true "Password length (5-100)"
// @Param special query bool true "Using special symbols (true or false)"
// @Success 200 {object} string
// @Failure 400 {object} dto.ErrorResponse
// @Router /gen [get]
func (h *PasswordHandler) Generate(w http.ResponseWriter, r *http.Request) {
	ps, err := getPasswordSettingsFromURL(r)
	if err != nil {
		logrus.Errorf(logs.FailedToGetPasswordSettings, err)
		sendErrorRespose(w, r, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	password, err := h.PasswordUseCase.Generate(ps)
	if err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
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
