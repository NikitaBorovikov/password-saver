package server

import (
	"context"
	"net/http"
	"password-saver/internal/infrastructure/api/handlers"
	"password-saver/internal/infrastructure/api/routes"
	"password-saver/internal/infrastructure/config"

	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(h handlers.Handlers, cfg *config.Config) *Server {
	router := routes.InitRoutes(h, cfg)
	srv := &Server{
		httpServer: &http.Server{
			Addr:        cfg.Http.Port,
			Handler:     router,
			ReadTimeout: cfg.Http.ReadTimeout,
			IdleTimeout: cfg.Http.IdleTimeout,
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
	logrus.Infof("Server started on port %s", s.httpServer.Addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
