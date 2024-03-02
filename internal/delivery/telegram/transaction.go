package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) NewTransaction(ctx *Context) {
	msg := tgbotapi.NewMessage(ctx.UID, "Please enter amount of coin to buy")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Hello", "some data")})
	_, err := h.bot.Send(msg)
	if err != nil {
		h.logger.Error(err)
	}
}
