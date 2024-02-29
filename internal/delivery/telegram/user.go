package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CreateUser(update *tgbotapi.Update) {
	uid := update.SentFrom().ID

	_, err := h.service.User.Create(uid)
	if err != nil {
		h.ResponseString(uid, h.startError(err))
		return
	}

	h.ResponseString(uid, h.startSuccess())
}

func (h *Handler) startError(err error) string {
	return fmt.Sprintf("Sorry, I can't create new user, reason: %v", err)
}

func (h *Handler) startSuccess() string {
	return "Congratulations! New user successfully created"
}
