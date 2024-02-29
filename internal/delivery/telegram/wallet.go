package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"math"
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
	return fmt.Sprintf("Sorry, I can't create new wallet, reason: %v", err)
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

	msg := tgbotapi.NewMessage(uid, "There is the list of your wallets")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(getKeyboardFromWallets(all)...)
	h.Response(msg)
}

func (h *Handler) allWalletsError(err error) string {
	return fmt.Sprintf("Sorry, I retrieve your wallets, reason: %v", err.Error())
}

func getKeyboardFromWallets(wallets []db.CryptoWallet) [][]tgbotapi.KeyboardButton {
	rows := math.Ceil(float64(len(wallets)) / float64(2))
	keyboard := make([][]tgbotapi.KeyboardButton, int(rows))

	var row int
	var col int

	for _, w := range wallets {
		keyboard[row] = append(keyboard[row], tgbotapi.NewKeyboardButton(fmt.Sprintf("/%s %s", walletMessage, w.Name)))

		if (col+1)%2 == 0 {
			col = 0
			row++
		} else {
			col++
		}
	}

	return keyboard
}

func (h *Handler) Wallet(update *tgbotapi.Update) {
	uid := update.SentFrom().ID
	symbols := strings.Split(update.Message.CommandArguments(), " ")
	if len(symbols) != 1 {
		h.ResponseString(uid, h.walletInvalidArguments())
		return
	}

	wallet, err := h.service.Wallet.Get(uid, symbols[0])
	if err != nil {
		h.ResponseString(uid, h.walletError(err))
		return
	}

	h.ResponseString(uid, h.walletSuccess(wallet))
}

func (h *Handler) walletInvalidArguments() string {
	return "Invalid command argument, expected: /wallet <coin_symbol>"
}

func (h *Handler) walletError(err error) string {
	return fmt.Sprintf("Sorry... I cant retrieve information about your wallet, reason: %v", err)
}

func (h *Handler) walletSuccess(w models.Wallet) string {
	return fmt.Sprintf(
		"Your wallet information:\n\t- Name: %s\n\t- Price: %f\n\t- Amount: %f\n\t- Balance: %f",
		w.Name, w.Price, w.Amount, w.Balance,
	)
}
