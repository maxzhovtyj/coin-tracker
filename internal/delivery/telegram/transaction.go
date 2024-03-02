package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"math"
	"strconv"
)

func (h *Handler) NewTransaction(ctx *Context) {
	all, err := h.service.Wallet.All(ctx.UID)
	if err != nil {
		ctx.ResponseString(fmt.Sprintf("Sorry, cannot get your wallet, %v", err))
		return
	}

	msg := tgbotapi.NewMessage(ctx.UID, "Please select wallet")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(getKeyboardFromWallets(all)...)

	_, err = ctx.ResponseWithError(msg)
	if err != nil {
		h.logger.Error(err)
		return
	}

	ctx.FSM.Update(ctx.UID, State{
		Command: newTransactionCommand,
		Step:    selectWalletStep,
	})
}

func getKeyboardFromWallets(wallets []db.CryptoWallet) [][]tgbotapi.KeyboardButton {
	rows := math.Ceil(float64(len(wallets)) / float64(2))
	keyboard := make([][]tgbotapi.KeyboardButton, int(rows))

	var row int
	var col int

	for _, w := range wallets {
		keyboard[row] = append(keyboard[row], tgbotapi.NewKeyboardButton(w.Name))

		if (col+1)%2 == 0 {
			col = 0
			row++
		} else {
			col++
		}
	}

	return keyboard
}

const (
	selectWalletStep = "select wallet"
	inputAmountStep  = "input amount"
	inputPriceStep   = "input price"
)

type NewTransactionData struct {
	Wallet models.Wallet
	Price  float64
	Amount float64
}

func (h *Handler) ResolveNewTransactionSteps(ctx *Context) {
	state := ctx.FSM.Get(ctx.UID)

	switch state.Step {
	case selectWalletStep:
		h.selectWalletStep(ctx)
	case inputAmountStep:
		h.inputTransactionAmountStep(ctx)
	case inputPriceStep:
		h.inputTransactionPriceStep(ctx)
	}
}

func (h *Handler) selectWalletStep(ctx *Context) {
	wallet, err := h.service.Wallet.Get(ctx.UID, h.resolveWalletName(ctx.Update.Message.Text))
	if err != nil {
		ctx.ResponseString(fmt.Sprintf("Sorry, I can't find this wallet, %v", err))
		return
	}

	ctx.ResponseString("Please enter amount (number)")

	ctx.FSM.Update(ctx.UID, State{
		Command: newTransactionCommand,
		Step:    inputAmountStep,
		Data: NewTransactionData{
			Wallet: wallet,
		},
	})
}

func (h *Handler) inputTransactionAmountStep(ctx *Context) {
	data, ok := ctx.FSM.Get(ctx.UID).Data.(NewTransactionData)
	if !ok {
		ctx.ResponseString("Error while processing command")
		ctx.FSM.Remove(ctx.UID)
		return
	}

	amount, err := strconv.ParseFloat(ctx.Update.Message.Text, 64)
	if err != nil {
		ctx.ResponseString("Invalid amount, expected number, try again")
		return
	}

	ctx.ResponseString("Input coin price")
	data.Amount = amount
	ctx.FSM.Update(ctx.UID, State{
		Command: newTransactionCommand,
		Step:    inputPriceStep,
		Data:    data,
	})
}

func (h *Handler) inputTransactionPriceStep(ctx *Context) {
	data, ok := ctx.FSM.Get(ctx.UID).Data.(NewTransactionData)
	if !ok {
		ctx.ResponseString("Error while processing command")
		ctx.FSM.Remove(ctx.UID)
		return
	}

	price, err := strconv.ParseFloat(ctx.Update.Message.Text, 64)
	if err != nil {
		ctx.ResponseString("Invalid price, expected number, try again")
		return
	}

	err = h.service.Wallet.NewTransaction(data.Wallet.Id, data.Amount, price)
	if err != nil {
		ctx.ResponseString(fmt.Sprintf("Failed to create transaction, %v", err))
		return
	}

	ctx.ResponseString("Transaction successfully saved")
	ctx.FSM.Remove(ctx.UID)
}
