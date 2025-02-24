package handlers

import "github.com/gorilla/sessions"

var (
	sessionKey   string
	sessionStore *sessions.CookieStore
)

func InitSessionStore(key string) {
	sessionKey = key
	sessionStore = sessions.NewCookieStore([]byte(sessionKey))
}
