package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) CreateUser(update *tgbotapi.Update) {
	uid := update.SentFrom().ID

	_, err := h.service.User.Create(uid)
	if err != nil {
		h.ResponseString(uid, h.StartError(err))
		return
	}

	h.ResponseString(uid, h.StartSuccess())
}

func (h *Handler) StartError(err error) string {
	return fmt.Sprintf("Sorry, I can't create new user, reason: %v", err.Error())
}

func (h *Handler) StartSuccess() string {
	return "Congratulations! New user successfully created"
}
