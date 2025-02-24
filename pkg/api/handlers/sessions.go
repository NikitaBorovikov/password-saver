package handlers

import "github.com/gorilla/sessions"

var (
	sessionKey   string
	sessionStore *sessions.CookieStore
)

func InitSession(key string) {
	sessionKey = key
	sessionStore = sessions.NewCookieStore([]byte(sessionKey))
}
