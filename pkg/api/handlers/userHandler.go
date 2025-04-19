package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"password-saver/pkg/api/session"
	"password-saver/pkg/dto"
	"password-saver/pkg/logs"
	"password-saver/pkg/usecases"

	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	UserUseCase *usecases.UserUseCase
	Session     *session.SessionManager
}

func newUserHandler(uc *usecases.UserUseCase, session *session.SessionManager) *UserHandler {
	return &UserHandler{
		UserUseCase: uc,
		Session:     session,
	}
}

// @Summary Register a new user
// @Description Creates a new user account with email and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.AuthRequest true "User registration data"
// @Success 201 {object} dto.GetUserInfoResponse
// @Failure 400,422 {object} dto.ErrorResponse
// @Router /auth/reg [post]
func (h *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {

	req, err := decodeAuthRequest(r)
	if err != nil {
		logrus.Errorf(logs.FailedToDecodeRequest, err)
		sendErrorRespose(w, r, http.StatusBadRequest, ErrDecodeRequest)
		return
	}

	userID, err := h.UserUseCase.Registration(req)
	if err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
		return
	}

	userInfo := dto.NewGetUserInfoResponse(userID, req.Email)

	sendOKResponse(w, r, http.StatusCreated, userInfo)

	logrus.Infof(logs.UserRegSuccessfully, userID)
}

// @Summary User authentication
// @Description Log in with user's username and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.AuthRequest true "User log in data"
// @Success 200 {integer} {userID}
// @Failure 400,401,500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) LogIn(w http.ResponseWriter, r *http.Request) {

	req, err := decodeAuthRequest(r)
	if err != nil {
		logrus.Errorf(logs.FailedToDecodeRequest, err)
		sendErrorRespose(w, r, http.StatusBadRequest, ErrDecodeRequest)
		return
	}

	user, err := h.UserUseCase.LogIn(req)
	if err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
		return
	}

	if err := setUserSession(w, r, h.Session, user.UserID); err != nil {
		logrus.Errorf("%s: %v", logs.FailedToGetSession, err)
		sendErrorRespose(w, r, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	sendOKResponse(w, r, http.StatusOK, user.UserID)

	logrus.Infof(logs.UserLoginedSuccessfully, user.UserID)
}

// @Summary Update user's data
// @Description Update user's password by user ID from context (an active session is required).
// @Tags User
// @Accept json
// @Produce json
// @Param input body dto.UpdateUserRequest true "Old and new user's passwords"
// @Success 200
// @Failure 400,422,500 {object} dto.ErrorResponse
// @Security SessionCookie
// @Router /profile/ [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	req, err := decodeUpdateRequest(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, ErrDecodeRequest)
		return
	}

	if err := h.UserUseCase.Update(req, userID); err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
		return
	}

	sendOKResponse(w, r, http.StatusOK, nil)

	logrus.Infof(logs.UserUpdatedSuccessfully, userID)
}

// @Summary Get user by ID
// @Description Get user's data by user ID from context (an active session is required).
// @Tags User
// @Produce json
// @Success 200 {object} dto.GetUserInfoResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security SessionCookie
// @Router /profile/ [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	user, err := h.UserUseCase.GetByID(userID)
	if err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
		return
	}

	userInfo := dto.NewGetUserInfoResponse(user.UserID, user.Email)

	sendOKResponse(w, r, http.StatusOK, userInfo)

	logrus.Info(logs.UserGivenByIDSuccessfully)
}

// @Summary Delete user's profile
// @Description delete user's info by user id from context (an active session is required).
// @Tags User
// @Produce json
// @Success 204
// @Failure 500 {object} dto.ErrorResponse
// @Security SessionCookie
// @Router /profile/ [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error(logs.FailedToGetUserIDFromCtx)
		sendErrorRespose(w, r, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	if err := h.UserUseCase.Delete(userID); err != nil {
		statusCode, apiErr := handleUsecaseErrors(err)
		sendErrorRespose(w, r, statusCode, apiErr)
		return
	}

	sendOKResponse(w, r, http.StatusNoContent, nil)

	logrus.Infof(logs.UserDeletedSuccessfully, userID)
}

// @Summary Log out of the system
// @Description Disables the user's session (an active session is required).
// @Tags User
// @Produce json
// @Success 204
// @Failure 500 {object} dto.ErrorResponse
// @Security SessionCookie
// @Router /profile/logout [post]
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	s := h.Session
	session, err := s.Store.Get(r, s.Name)
	if err != nil || session == nil {
		logrus.Errorf("%s: %v", logs.FailedToGetSession, err)
		sendErrorRespose(w, r, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	cleanSessionInfo(session)

	if err := session.Save(r, w); err != nil {
		logrus.Errorf("%s: %v", logs.FailedToSaveSession, err)
		sendErrorRespose(w, r, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	sendOKResponse(w, r, http.StatusNoContent, nil)

}

func decodeAuthRequest(r *http.Request) (*dto.AuthRequest, error) {
	var req dto.AuthRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeUpdateRequest(r *http.Request) (*dto.UpdateUserRequest, error) {
	var req dto.UpdateUserRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func setUserSession(w http.ResponseWriter, r *http.Request, s *session.SessionManager, userID int64) error {

	session, err := s.Store.Get(r, s.Name)
	if err != nil || session == nil {
		return fmt.Errorf(logs.FailedToGetSessionKey, err)
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return fmt.Errorf("failed to generate random bytes array: %v", err)
	}

	setSessionValues(session, sessionID, userID)

	err = saveSession(session, r, w)
	return err
}

func generateSessionID() (string, error) {
	byteArr := make([]byte, 32)

	_, err := rand.Read(byteArr)
	if err != nil {
		return "", err
	}

	sessionID := hex.EncodeToString(byteArr)
	return sessionID, nil
}

func setSessionValues(session *sessions.Session, sessionID string, userID int64) {
	session.Values[sessionAuthenticated] = true
	session.Values[sessionIDKey] = sessionID
	session.Values[sessionUserIDKey] = userID
}

func saveSession(session *sessions.Session, r *http.Request, w http.ResponseWriter) error {
	err := session.Save(r, w)
	return err
}

func cleanSessionInfo(session *sessions.Session) {
	session.Values[sessionAuthenticated] = false
	session.Options.MaxAge = -1 // delete session
}
