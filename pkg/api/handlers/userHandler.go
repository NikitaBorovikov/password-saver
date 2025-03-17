package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"password-saver/pkg/api/session"
	"password-saver/pkg/dto"
	apperrors "password-saver/pkg/errors"
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

func (h *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {

	req, err := decodeAuthRequest(r)
	if err != nil {
		logrus.Errorf("failed to decode request: %v", err)
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrDecodeRequest)
		return
	}

	userID, err := h.UserUseCase.Registration(req)
	if err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	userInfo := dto.NewGetUserInfoResponse(userID, req.Email)

	sendOKResponse(w, r, http.StatusCreated, userInfo)

	logrus.Infof("user was registated successfully with id = %d", userID)
}

func (h *UserHandler) LogIn(w http.ResponseWriter, r *http.Request) {

	req, err := decodeAuthRequest(r)
	if err != nil {
		logrus.Errorf("failed to decode request: %v", err)
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrDecodeRequest)
		return
	}

	user, err := h.UserUseCase.LogIn(req)
	if err != nil {
		sendErrorRespose(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := setUserSession(w, r, h.Session, user.UserID); err != nil {
		logrus.Errorf("failed to set session: %v", err)
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	sendOKResponse(w, r, http.StatusOK, user.UserID)

	logrus.Infof("user {id = %d} was lodin successfully", user.UserID)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error("failed to get userID from context")
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	req, err := decodeUpdateRequest(r)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, apperrors.ErrDecodeRequest)
		return
	}

	if err := h.UserUseCase.Update(req, userID); err != nil {
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	sendOKResponse(w, r, http.StatusOK, "user data was updated")

	logrus.Infof("user {id = %d} was updated successfully", userID)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error("failed to get userID from context")
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	user, err := h.UserUseCase.GetByID(userID)
	if err != nil {
		sendErrorRespose(w, r, http.StatusInternalServerError, err)
		return
	}

	userInfo := dto.NewGetUserInfoResponse(user.UserID, user.Email)

	sendOKResponse(w, r, http.StatusOK, userInfo)

	logrus.Info("successfull getting user data by id")
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		logrus.Error("failed to get userID from context")
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	if err := h.UserUseCase.Delete(userID); err != nil {
		sendErrorRespose(w, r, http.StatusInternalServerError, err)
		return
	}

	sendOKResponse(w, r, http.StatusNoContent, nil)

	logrus.Infof("user {id = %d} was deleted successfully", userID)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	s := h.Session
	session, err := s.Store.Get(r, s.Name)
	if err != nil || session == nil {
		logrus.Errorf("failed to get user session: %v", err)
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	cleanSessionInfo(session)

	if err := session.Save(r, w); err != nil {
		logrus.Errorf("failed to save user session: %v", err)
		sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
		return
	}

	sendOKResponse(w, r, http.StatusNoContent, "logout is done")

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
		return fmt.Errorf("failed to get sessionKey: %v", err)
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
