package telegram

import (
	"password-saver/pkg/infrastructure/telegram/handlers"
	"password-saver/pkg/usecases"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	handlers *handlers.Handlers
}

func NewBot(bot *tgbotapi.BotAPI, uc *usecases.UseCases) *Bot {
	return &Bot{
		bot:      bot,
		handlers: handlers.InitHandlers(uc),
	}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	b.handleUpdates(updates)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				b.commandHandler(update)
			} else {

			}
		}
	}
}

func (b *Bot) commandHandler(update tgbotapi.Update) error {
	switch update.Message.Command() {
	case "start":
		return nil
	case "newPassword":
		return b.handlers.PasswordHandler.Generate(update, b.bot)
	default:
		return nil
	}
}
