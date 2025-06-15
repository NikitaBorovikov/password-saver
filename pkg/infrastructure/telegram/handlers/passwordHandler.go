package handlers

import (
	"password-saver/pkg/infrastructure/dto"
	"password-saver/pkg/usecases"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PasswordHandler struct {
	PasswordUseCase *usecases.PasswordUseCase
}

func newPasswordHandler(uc *usecases.PasswordUseCase) *PasswordHandler {
	return &PasswordHandler{
		PasswordUseCase: uc,
	}
}

func (h *PasswordHandler) Generate(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	req := &dto.GeneratePasswordRequest{
		Length:            15,
		UseSpecialSymbols: true,
	}
	password, err := h.PasswordUseCase.Generate(req)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error")
		bot.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, password)
	bot.Send(msg)
	return nil
}
