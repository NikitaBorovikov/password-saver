package main

import (
	"context"
	"os"
	"os/signal"
	"password-saver/internal/infrastructure/api/handlers"
	"password-saver/internal/infrastructure/api/server"
	"password-saver/internal/infrastructure/api/session"
	"password-saver/internal/infrastructure/config"
	"password-saver/internal/infrastructure/db"
	"password-saver/internal/infrastructure/logs"
	"password-saver/internal/infrastructure/repository"
	"password-saver/internal/usecases"
	"syscall"

	"github.com/sirupsen/logrus"
)

// @title Password-Saver API
// @version 1.0
// @description REST API for secure storage and generation of new passwords

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey SessionCookie
// @in cookie
// @name password_saver_auth_session
// @description Some endpoints require an active session (the session is stored for 12 hours).
func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.Init()
	if err != nil {
		logrus.Fatalf(logs.FailedToInitConfig, err)
	}

	db, err := db.ConnAndPing(cfg.Postgres)
	if err != nil {
		logrus.Fatalf(logs.FailedToConnectDB, err)
	}
	defer db.Close()

	session := session.NewSessionManager(cfg.Http.SessionKey, cfg.Http.SessionName)

	repo := repository.InitRepository(db)
	usecases := usecases.InitUseCases(repo.UserRepository, repo.PasswordRepository, repo.SystemRepository, &cfg.EncryptKeys)
	handlers := handlers.InitHandlers(usecases, session)

	srv := server.NewServer(*handlers, cfg)
	srv.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf(logs.FailedShutDownServer, err)
	}
}
