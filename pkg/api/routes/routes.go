package routes

import (
	"password-saver/pkg/api/handlers"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(h handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/auth", func(r chi.Router) {
		authRoutes(r, h)
	})
	return r
}

func authRoutes(r chi.Router, h handlers.Handlers) {
	r.Post("/reg", h.UserHandler.Registration)
	r.Post("/login", h.UserHandler.LogIn)
}
