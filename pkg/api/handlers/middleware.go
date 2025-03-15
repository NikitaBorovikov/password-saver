package handlers

import (
	"context"
	"net/http"
	"password-saver/pkg/api/session"
	apperrors "password-saver/pkg/errors"

	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

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
				logrus.Errorf("failed to get session: %v", err)
				sendErrorRespose(w, r, http.StatusInternalServerError, apperrors.ErrServerInternal)
				return
			}

			auth, ok := session.Values[sessionAuthenticated].(bool)
			if !ok || !auth {
				logrus.Error("user is unauthenticated")
				sendErrorRespose(w, r, http.StatusUnauthorized, apperrors.ErrNotAuthenticated)
				return
			}

			userID, ok := session.Values[sessionUserIDKey].(int64)
			if !ok {
				logrus.Error("userID not in session")
				sendErrorRespose(w, r, http.StatusUnauthorized, apperrors.ErrNotAuthenticated)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDCtx, userID)
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

func LoggingMiddleWare() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			logrus.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
			}).Info("Incoming request")

			next.ServeHTTP(w, r)

		})
	}
}
