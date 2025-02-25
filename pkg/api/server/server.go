package server

import (
	"context"
	"net/http"
	"password-saver/pkg/api/handlers"
	"password-saver/pkg/api/routes"
	"password-saver/pkg/config"

	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(h handlers.Handlers, cfg *config.Http) *Server {
	router := routes.InitRoutes(h)
	handlers.InitSessionStore(cfg.SessionKey, "auth") // where to save a session name?

	srv := &Server{
		httpServer: &http.Server{
			Addr:    cfg.Port,
			Handler: router,
		},
	}

	return srv
}

func (s *Server) Run() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			logrus.Fatal(err)
		}
	}()
	logrus.Printf("Server started on port %s", s.httpServer.Addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
