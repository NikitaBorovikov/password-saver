package routes

import (
	"password-saver/pkg/api/handlers"

	_ "password-saver/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRoutes(h handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware stack
	r.Use(
		handlers.CORSMiddleware(),
		handlers.LoggingMiddleware(),
	)
	authMiddleware := handlers.AuthMiddleware(h.UserHandler.Session)

	// Auth routes with strict rate limiting
	r.Route("/auth", func(r chi.Router) {
		r.Use(handlers.RateLimiterMiddleware(5))
		authRoutes(r, *h.UserHandler)
	})

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(handlers.RateLimiterMiddleware(30))

		// Profile routes
		r.Route("/profile", func(r chi.Router) {
			profileRoutes(r, *h.UserHandler)
		})

		// Password routes
		r.Route("/passwords", func(r chi.Router) {
			passwordRoutes(r, *h.PasswordHandler)
		})

	})

	// Open routes
	r.Group(func(r chi.Router) {
		r.Use(handlers.RateLimiterMiddleware(30))
		openRouters(r, h)
	})

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))

	return r
}

func authRoutes(r chi.Router, h handlers.UserHandler) {
	r.Post("/reg", h.Registration)
	r.Post("/login", h.LogIn)
}

func profileRoutes(r chi.Router, h handlers.UserHandler) {
	r.Get("/", h.GetByID)
	r.Delete("/", h.Delete)
	r.Put("/", h.Update)
	r.Post("/logout", h.Logout)
}

func passwordRoutes(r chi.Router, h handlers.PasswordHandler) {
	r.Post("/", h.Save)
	r.Get("/", h.GetAll)

	r.Route("/{passwordID}", func(r chi.Router) {
		r.Get("/", h.GetByID)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})
}

func openRouters(r chi.Router, h handlers.Handlers) {
	r.Get("/gen", h.PasswordHandler.Generate)
}
