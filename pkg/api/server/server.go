package server

import (
	"net/http"
	"password-saver/pkg/api/handlers"
	"password-saver/pkg/api/routes"
	"password-saver/pkg/config"

	"github.com/sirupsen/logrus"
)

type server struct {
	httpServer *http.Server
	handlers   handlers.Handlers
	cfg        config.Config
}

func NewServer(handlers handlers.Handlers, cfg *config.Config) *server {

	router := routes.InitRoutes(handlers)

	httpServer := &http.Server{
		Addr:    cfg.Http.Port,
		Handler: router,
	}

	srv := &server{
		httpServer: httpServer,
		handlers:   handlers,
		cfg:        *cfg,
	}

	return srv
}

func (s *server) Start() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			logrus.Fatalf("Server failed to start: %s", err.Error())
		}
	}()
	logrus.Printf("Server started on port %s", s.cfg.Http.Port)
}
