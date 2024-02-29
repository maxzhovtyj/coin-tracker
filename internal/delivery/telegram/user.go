package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CreateUser(update *tgbotapi.Update) tgbotapi.MessageConfig {
	user, err := h.service.User.Create(update.SentFrom().ID)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(`User with id %d successfully created`, user.TelegramID))
}
