package main

import (
	"password-saver/pkg/config"
	"password-saver/pkg/db"
	"password-saver/pkg/handlers/api"
	"password-saver/pkg/repository"
	"password-saver/pkg/usecases"

	"github.com/sirupsen/logrus"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := db.ConnAndPing(cfg.Postgres)
	if err != nil {
		logrus.Fatal(err)
	}

	repository := repository.InitRepository(db)
	usecases := usecases.InitUseCases(repository)
	handlers := api.InitHandlers(usecases)

	_ = handlers

}
