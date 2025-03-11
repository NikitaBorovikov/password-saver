package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"password-saver/pkg/api/session"

	"github.com/go-chi/cors"
)

type contextKey string

const UserIDKey contextKey = "userID"

const (
	sessionAuthenticated = "authenticated"
	sessionIDKey         = "sessionID"
	sessionUserIDKey     = "userID"
)

func AuthMiddleware(sm *session.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sm.Store.Get(r, sm.Name)
			if err != nil || session == nil {
				log.Printf("Failed to get session: %v", err)
				sendErrorRespose(w, r, http.StatusInternalServerError, err)
				return
			}

			auth, ok := session.Values[sessionAuthenticated].(bool)
			if !ok || !auth {
				sendErrorRespose(w, r, http.StatusUnauthorized, fmt.Errorf("not authenticated"))
				return
			}

			userID, ok := session.Values[sessionUserIDKey].(int64)
			if !ok {
				err := fmt.Errorf("user ID not found or invalid")
				sendErrorRespose(w, r, http.StatusUnauthorized, err)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func CORSMiddleware() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})
}
