package handlers

import "github.com/gorilla/sessions"

var (
	sessionName  string
	sessionStore *sessions.CookieStore
)

func InitSessionStore(key, name string) {
	sessionName = name
	sessionStore = sessions.NewCookieStore([]byte(key))
}
