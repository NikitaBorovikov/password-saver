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
	"password-saver/pkg/repository"
	"password-saver/pkg/usecases"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		logrus.Fatalf("config failed to init: %v", err)
	}

	db, err := db.ConnAndPing(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("db failed to connect: %v", err)
	}
	defer db.Close()

	session := session.NewSessionManager(cfg.Http.SessionKey, "auth")

	repository := repository.InitRepository(db)
	usecases := usecases.InitUseCases(repository, &cfg.EncryptKeys)
	handlers := handlers.InitHandlers(usecases, session)

	srv := server.NewServer(*handlers, &cfg.Http)
	srv.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

}
