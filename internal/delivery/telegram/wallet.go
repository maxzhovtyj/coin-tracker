package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (h *Handler) NewWallet(update *tgbotapi.Update) {
	uid := update.SentFrom().ID
	symbols := strings.Split(update.Message.CommandArguments(), " ")
	if len(symbols) != 1 {
		h.Response(uid, h.InvalidNewWalletArguments())
		return
	}

	err := h.service.Wallet.Create(uid, symbols[0])
	if err != nil {
		h.Response(uid, h.NewWalletError(err))
		return
	}

	h.Response(uid, h.NewWalletSuccess())
}

func (h *Handler) InvalidNewWalletArguments() string {
	return "Invalid command argument, expected: /newWallet <coin_symbol>"
}

func (h *Handler) NewWalletError(err error) string {
	return fmt.Sprintf("Sorry, I can't create new wallet, reason: %v", err.Error())
}

func (h *Handler) NewWalletSuccess() string {
	return "Congratulations! New wallet successfully created"
}
