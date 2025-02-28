package handlers

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessionStore.Get(r, sessionName)
		if err != nil {
			sendErrorRespose(w, r, http.StatusInternalServerError, err)
			return
		}

		auth, ok := session.Values[sessionAuthenticated].(bool)
		if !ok || !auth {
			sendErrorRespose(w, r, http.StatusForbidden, err)
			return
		}

		userID := session.Values[sessionUserIDKey].(int64)
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
