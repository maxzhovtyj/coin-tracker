package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Response(chatID int64, msg string) {
	_, err := h.bot.Send(tgbotapi.NewMessage(chatID, msg))
	if err != nil {
		h.logger.Errorf("failed to send message: %v", err)
		return
	}
}

func (h *Handler) UnknownCommand() string {
	return "Sorry, I don't what command it is"
}
