package main

import (
	"password-saver/pkg/config"
	"password-saver/pkg/db"
	"password-saver/pkg/infrastructure/repository"
	"password-saver/pkg/infrastructure/telegram"
	"password-saver/pkg/logs"
	"password-saver/pkg/usecases"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
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

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		logrus.Fatal("failed to init tg bot")
	}
	bot.Debug = true

	repo := repository.InitRepository(db)
	usecases := usecases.InitUseCases(repo.UserRepository, repo.PasswordRepository, repo.SystemRepository, &cfg.EncryptKeys)
	tgBot := telegram.NewBot(bot, usecases)

	tgBot.Start()

}
