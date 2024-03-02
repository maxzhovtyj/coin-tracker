package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"math"
)

func (h *Handler) NewWallet(ctx *Context) {
	symbols := ctx.CommandArgs()
	if len(symbols) != 1 {
		ctx.ResponseString(h.invalidNewWalletArguments())
		return
	}

	walletName := h.resolveWalletName(symbols[0])

	err := h.service.Wallet.Create(ctx.UID, walletName)
	if err != nil {
		ctx.ResponseString(h.newWalletError(err))
		return
	}

	ctx.ResponseString(h.newWalletSuccess())
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

func (h *Handler) Wallets(ctx *Context) {
	all, err := h.service.Wallet.All(ctx.UID)
	if err != nil {
		ctx.ResponseString(h.allWalletsError(err))
		return
	}

	if len(all) == 0 {
		ctx.ResponseString("Sorry, you don't have any wallets yet, use /newWallet command")
		return
	}

	msg := tgbotapi.NewMessage(ctx.UID, "There is the list of your wallets")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(getKeyboardFromWallets(all)...)
	ctx.Response(msg)
}

func (h *Handler) allWalletsError(err error) string {
	return fmt.Sprintf("Sorry, I retrieve your wallets, reason: %v", err.Error())
}

func getKeyboardFromWallets(wallets []db.CryptoWallet) [][]tgbotapi.InlineKeyboardButton {
	rows := math.Ceil(float64(len(wallets)) / float64(2))
	keyboard := make([][]tgbotapi.InlineKeyboardButton, int(rows))

	var row int
	var col int

	for _, w := range wallets {
		cbData := fmt.Sprintf("%s=%s", walletCallback, w.Name)
		keyboard[row] = append(keyboard[row], tgbotapi.NewInlineKeyboardButtonData(w.Name, cbData))

		if (col+1)%2 == 0 {
			col = 0
			row++
		} else {
			col++
		}
	}

	return keyboard
}

func (h *Handler) Wallet(ctx *Context) {
	walletName := h.resolveWalletName(ctx.CallbackDataValue)

	wallet, err := h.service.Wallet.Get(ctx.UID, walletName)
	if err != nil {
		ctx.ResponseString(h.walletError(err))
		return
	}

	ctx.ResponseString(h.walletSuccess(wallet))
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
