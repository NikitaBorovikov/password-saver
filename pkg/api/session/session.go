package session

import "github.com/gorilla/sessions"

type SessionManager struct {
	Store *sessions.CookieStore
	Name  string
}

func NewSessionManager(key, name string) *SessionManager {
	store := sessions.NewCookieStore([]byte(key))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 12, // 12 hours
		HttpOnly: true,
	}

	return &SessionManager{
		Store: store,
		Name:  name,
	}
}
