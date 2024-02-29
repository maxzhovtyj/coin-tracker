package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (h *Handler) NewWallet(update *tgbotapi.Update) tgbotapi.MessageConfig {
	err := h.service.Wallet.Create(update.SentFrom().ID, "BTCUSDT")
	if err != nil {
		return tgbotapi.NewMessage(update.SentFrom().ID, err.Error())
	}

	return tgbotapi.NewMessage(update.SentFrom().ID, "New wallet successfully created")
}
