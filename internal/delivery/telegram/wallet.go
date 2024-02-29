package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"strings"
)

func (h *Handler) NewWallet(update *tgbotapi.Update) {
	uid := update.SentFrom().ID
	symbols := strings.Split(update.Message.CommandArguments(), " ")
	if len(symbols) != 1 {
		h.ResponseString(uid, h.invalidNewWalletArguments())
		return
	}

	err := h.service.Wallet.Create(uid, symbols[0])
	if err != nil {
		h.ResponseString(uid, h.newWalletError(err))
		return
	}

	h.ResponseString(uid, h.newWalletSuccess())
}

func (h *Handler) invalidNewWalletArguments() string {
	return "Invalid command argument, expected: /newWallet <coin_symbol>"
}

func (h *Handler) newWalletError(err error) string {
	return fmt.Sprintf("Sorry, I can't create new wallet, reason: %v", err.Error())
}

func (h *Handler) newWalletSuccess() string {
	return "Congratulations! New wallet successfully created"
}

func (h *Handler) Wallets(update *tgbotapi.Update) {
	uid := update.SentFrom().ID

	all, err := h.service.Wallet.All(uid)
	if err != nil {
		h.ResponseString(uid, h.allWalletsError(err))
		return
	}

	msg := tgbotapi.NewMessage(uid, "There are list of your wallets")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(h.allWalletsSuccess(all)...)
	h.Response(msg)
}

func (h *Handler) allWalletsError(err error) string {
	return fmt.Sprintf("Sorry, I retrieve your wallets, reason: %v", err.Error())
}

func (h *Handler) allWalletsSuccess(wallets []db.CryptoWallet) [][]tgbotapi.KeyboardButton {
	keyboard := make([][]tgbotapi.KeyboardButton, len(wallets)/3)

	var row int
	var col int

	for _, w := range wallets {
		keyboard[row] = append(keyboard[row], tgbotapi.NewKeyboardButton(w.Name))

		col++
		if col+1%3 == 0 {
			col = 0
			row++
		}
	}

	return keyboard
}

func (h *Handler) NewWalletRecord(update *tgbotapi.Update) {

}
