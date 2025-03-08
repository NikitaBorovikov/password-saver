package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"password-saver/pkg/api/session"
	"password-saver/pkg/dto"
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

	req, err := decodeRegRequest(r)
	if err != nil {
		logrus.Error(err)
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	userID, err := h.UserUseCase.Registration(req)
	if err != nil {
		logrus.Error(err)
		sendErrorRespose(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	userInfo := dto.NewGetUserInfoResponse(userID, req.Email)

	sendOKResponse(w, r, userInfo)
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

	if err := setUserSession(w, r, h.Session, user.UserID); err != nil {
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

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		err := fmt.Errorf("no userID in context")
		sendErrorRespose(w, r, http.StatusUnauthorized, err)
		return
	}

	user, err := h.UserUseCase.GetByID(userID)
	if err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	userInfo := dto.NewGetUserInfoResponse(user.UserID, user.Email)

	sendOKResponse(w, r, userInfo)

}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		err := fmt.Errorf("no userID in context")
		sendErrorRespose(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.UserUseCase.Delete(userID); err != nil {
		sendErrorRespose(w, r, http.StatusBadRequest, err)
		return
	}

	sendOKResponse(w, r, nil)
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
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("failed to save session: %v", err)
	}
	return nil
}
