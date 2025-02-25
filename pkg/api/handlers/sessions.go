package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"password-saver/pkg/model"

	"github.com/gorilla/sessions"
)

const (
	sessionAuthenticated = "authenticated"
	sessionIDKey         = "sessionID"
	userIDKey            = "userID"
)

var (
	sessionName  string
	sessionStore *sessions.CookieStore
)

func InitSessionStore(key, name string) {
	sessionName = name
	sessionStore = sessions.NewCookieStore([]byte(key))
}

func setUserSession(w http.ResponseWriter, r *http.Request, u *model.User) error {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		return fmt.Errorf("failed to get sessionKey: %v", err)
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return fmt.Errorf("failed to generate random bytes array: %v", err)
	}

	setSessionValues(session, sessionID, u.UserID)

	setSessionOptions(session)

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
	session.Values[userIDKey] = userID
}

func setSessionOptions(session *sessions.Session) {
	session.Options = &sessions.Options{
		MaxAge:   3600 * 12,
		HttpOnly: true,
	}
}

func saveSession(session *sessions.Session, r *http.Request, w http.ResponseWriter) error {
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("failed to save session: %v", err)
	}
	return nil
}
