package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	sessionAuthenticated = "authenticated"
	sessionIDKey         = "sessionID"
	sessionUserIDKey     = "userID" // maybe rename?
)

var (
	sessionName  string
	sessionStore *sessions.CookieStore
)

func InitSessionStore(key, name string) {
	sessionName = name
	sessionStore = sessions.NewCookieStore([]byte(key))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 12, // 12 часов
		HttpOnly: true,      // Защита от XSS
	}
}

func setUserSession(w http.ResponseWriter, r *http.Request, userID int64) error {
	session, err := sessionStore.Get(r, sessionName)
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
