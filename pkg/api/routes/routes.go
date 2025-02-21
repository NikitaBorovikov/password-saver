package routes

import (
	"password-saver/pkg/api/handlers"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(h handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()

	return r
}
