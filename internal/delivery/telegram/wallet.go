package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"math"
	"strconv"
)

func (h *Handler) NewWallet(ctx *Context) {
	ctx.ResponseString("Please enter wallet name (ex. BTCUSDT, ETHUSDT)")

	ctx.FSM.Update(ctx.UID, State{
		Caller: ctx.CallbackName,
		Step:   "wallet input",
	})
}

func (h *Handler) ResolveNewWalletSteps(ctx *Context) {
	err := h.service.Wallet.Create(ctx.UID, h.resolveWalletName(ctx.Update.Message.Text))
	if err != nil {
		ctx.ResponseString(h.newWalletError(err))
		return
	}

	ctx.ResponseString(h.newWalletSuccess())
}

func (h *Handler) DeleteWallet(ctx *Context) {
	walletID, err := strconv.ParseInt(ctx.CallbackDataValue, 10, 64)
	if err != nil {
		ctx.ResponseString("Invalid wallet id")
		return
	}

	ctx.ResponseString("You are going to delete and all it's transactions, in order to do this enter full wallet name (ex. BTCUSDT, ETHUSDT)")

	ctx.FSM.Update(ctx.UID, State{
		Caller: ctx.CallbackName,
		Data:   walletID,
		Step:   "wallet delete",
	})
}

func (h *Handler) ResolveDeleteWalletStep(ctx *Context) {
	state := ctx.FSM.Get(ctx.UID)
	defer ctx.FSM.Remove(ctx.UID)

	walletID, ok := state.Data.(int64)
	if !ok {
		ctx.ResponseString("Invalid wallet id")
		return
	}

	walletName := ctx.Update.Message.Text
	wallet, err := h.service.Wallet.Get(ctx.UID, walletID)
	if err != nil {
		ctx.ResponseString("Can't get your wallet, " + err.Error())
		return
	}

	if wallet.Name != walletName {
		ctx.ResponseString("Wallet name mismatch, try again")
		return
	}

	err = h.service.Wallet.Delete(ctx.UID, walletID)
	if err != nil {
		ctx.ResponseString("Can't delete your wallet, " + err.Error())
		return
	}

	ctx.ResponseString("Your wallet successfully deleted")
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
		ctx.ResponseString("Sorry, you don't have any wallets yet")
		return
	}

	msg := tgbotapi.NewMessage(ctx.UID, "There is the list of your wallets")
	newBtn := tgbotapi.NewInlineKeyboardButtonData("New Wallet", walletNewCallback+"=")
	btns := [][]tgbotapi.InlineKeyboardButton{{newBtn}}
	btns = append(btns, getWalletCallbackKeyboard(all)...)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(btns...)
	ctx.Response(msg)
}

func (h *Handler) allWalletsError(err error) string {
	return fmt.Sprintf("Sorry, I retrieve your wallets, reason: %v", err.Error())
}

func getWalletCallbackKeyboard(wallets []db.CryptoWallet) [][]tgbotapi.InlineKeyboardButton {
	rows := math.Ceil(float64(len(wallets)) / float64(2))
	keyboard := make([][]tgbotapi.InlineKeyboardButton, int(rows))

	var row int
	var col int

	for _, w := range wallets {
		cbData := fmt.Sprintf("%s=%d", walletCallback, w.ID)
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
	walletID, err := strconv.ParseInt(ctx.CallbackDataValue, 10, 64)
	if err != nil {
		ctx.ResponseString("Invalid wallet id")
		return
	}

	wallet, err := h.service.Wallet.Get(ctx.UID, walletID)
	if err != nil {
		ctx.ResponseString(h.walletError(err))
		return
	}

	msg := tgbotapi.NewMessage(ctx.UID, h.walletSuccess(wallet))

	buyBtn := tgbotapi.NewInlineKeyboardButtonData("Buy", fmt.Sprintf("%s=%d", walletBuyCallback, walletID))
	sellBtn := tgbotapi.NewInlineKeyboardButtonData("Sell", fmt.Sprintf("%s=%d", walletSellCallback, walletID))
	transactionsBtn := tgbotapi.NewInlineKeyboardButtonData("Transactions", fmt.Sprintf("%s=%d", walletTransactionsCallback, walletID))
	deleteBtn := tgbotapi.NewInlineKeyboardButtonData("Delete", fmt.Sprintf("%s=%d", walletDeleteCallback, walletID))

	markup := [][]tgbotapi.InlineKeyboardButton{
		{transactionsBtn},
		{buyBtn, sellBtn},
		{deleteBtn},
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(markup...)

	ctx.Response(msg)
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
