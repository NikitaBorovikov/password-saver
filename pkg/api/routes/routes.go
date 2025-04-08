package routes

import (
	"password-saver/pkg/api/handlers"

	_ "password-saver/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRoutes(h handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(handlers.CORSMiddleware())
	r.Use(handlers.LoggingMiddleWare())

	authMiddleware := handlers.AuthMiddleware(h.UserHandler.Session)

	r.Route("/auth", func(r chi.Router) {
		r.Use(handlers.RateLimiterMeddleWare(5)) //5 requests per minute
		authRoutes(r, *h.UserHandler)
	})

	r.Route("/profile", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(handlers.RateLimiterMeddleWare(30))
		profileRoutes(r, *h.UserHandler)
	})

	r.Route("/passwords", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(handlers.RateLimiterMeddleWare(50))
		passwordRoutes(r, *h.PasswordHandler)
	})

	//swagger
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
	r.Get("/gen", h.Generate)
	r.Post("/", h.Save)
	r.Get("/", h.GetAll)

	r.Route("/{passwordID}", func(r chi.Router) {
		r.Get("/", h.GetByID)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})
}
