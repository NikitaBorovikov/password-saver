package main

import (
	"context"
	"os"
	"os/signal"
	"password-saver/pkg/api/handlers"
	"password-saver/pkg/api/server"
	"password-saver/pkg/api/session"
	"password-saver/pkg/config"
	"password-saver/pkg/db"
	"password-saver/pkg/logs"
	"password-saver/pkg/repository"
	"password-saver/pkg/usecases"
	"syscall"

	"github.com/sirupsen/logrus"
)

const (
	sessionName = "password_saver_auth_session"
)

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

	session := session.NewSessionManager(cfg.Http.SessionKey, sessionName)

	repository := repository.InitRepository(db)
	usecases := usecases.InitUseCases(repository, &cfg.EncryptKeys)
	handlers := handlers.InitHandlers(usecases, session)

	srv := server.NewServer(*handlers, &cfg.Http)
	srv.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf(logs.FailedShutDownServer, err)
	}
}
