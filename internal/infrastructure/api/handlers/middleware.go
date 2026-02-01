package handlers

import (
	"context"
	"net/http"
	"password-saver/internal/infrastructure/api/session"
	"password-saver/internal/infrastructure/logs"
	"time"

	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
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
				logrus.Errorf("%s: %v", logs.FailedToGetSession, err)
				sendErrorRespose(w, r, http.StatusInternalServerError, ErrInternalServer)
				return
			}

			auth, ok := session.Values[sessionAuthenticated].(bool)
			if !ok || !auth {
				logrus.Error(logs.UnauthenticatedUser)
				sendErrorRespose(w, r, http.StatusUnauthorized, ErrNotAuthenticated)
				return
			}

			userID, ok := session.Values[sessionUserIDKey].(int64)
			if !ok {
				logrus.Error(logs.FailedToGetUserIDFromSession)
				sendErrorRespose(w, r, http.StatusUnauthorized, ErrNotAuthenticated)
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

// Limits the number of requests
func RateLimiterMiddleware(requestsPerMin int) func(http.Handler) http.Handler {
	return httprate.Limit(
		requestsPerMin, // request amount
		time.Minute,    // interval
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			return r.RemoteAddr, nil // key is IP
		}),
	)
}

func LoggingMiddleware() func(http.Handler) http.Handler {
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
