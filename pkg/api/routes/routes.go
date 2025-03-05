package routes

import (
	"password-saver/pkg/api/handlers"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(h handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(handlers.CORSMiddleware())

	authMiddleware := handlers.AuthMiddleware(h.UserHandler.Session)

	r.Route("/auth", func(r chi.Router) {
		authRoutes(r, *h.UserHandler)
	})
	r.Route("/profile", func(r chi.Router) {
		r.Use(authMiddleware)
		profileRoutes(r, *h.UserHandler)
	})

	r.Route("/passwords", func(r chi.Router) {
		r.Use(authMiddleware)
		passwordRoutes(r, *h.PasswordHandler)
	})
	return r
}

func authRoutes(r chi.Router, h handlers.UserHandler) {
	r.Post("/reg", h.Registration)
	r.Post("/login", h.LogIn)
}

func profileRoutes(r chi.Router, h handlers.UserHandler) {
	r.Delete("/", h.Delete)
	r.Route("/{userID}", func(r chi.Router) {
		r.Put("/", h.Update)
	})
	//logout
}

func passwordRoutes(r chi.Router, h handlers.PasswordHandler) {
	r.Post("/", h.Save)
	// r.GetAll("/", h.GetAll)
	// r.Route("/{passwordID}, func(r chi.Rputer){
	// 		r.Get("/", h.GetByID)
	//		r.Put("/", h.Update)
	// }")
}
