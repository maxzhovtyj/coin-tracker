package telegram

import (
	"fmt"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	"github.com/maxzhovtyj/coin-tracker/internal/service"
	"math"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) NewTransaction(ctx *Context) {
	walletID, err := strconv.ParseInt(ctx.CallbackDataValue, 10, 64)
	if err != nil {
		ctx.ResponseString("Invalid wallet id")
		return
	}

	wallet, err := h.service.Wallet.Get(ctx.UID, walletID)
	if err != nil {
		ctx.ResponseString("Cant find wallet")
		return
	}

	data := NewTransactionData{
		Wallet: wallet,
	}

	ctx.ResponseString("Input coin amount")

	ctx.FSM.Update(ctx.UID, State{
		Caller: ctx.CallbackName,
		Step:   inputAmountStep,
		Data:   data,
	})
}

const (
	inputAmountStep = "input amount"
	inputPriceStep  = "input price"
)

type NewTransactionData struct {
	Wallet models.Wallet
	Price  float64
	Amount float64
}

func (h *Handler) ResolveNewTransactionSteps(ctx *Context) {
	state := ctx.FSM.Get(ctx.UID)

	switch state.Step {
	case inputAmountStep:
		h.inputTransactionAmountStep(ctx)
	case inputPriceStep:
		h.inputTransactionPriceStep(ctx)
	}
}

func (h *Handler) inputTransactionAmountStep(ctx *Context) {
	state := ctx.FSM.Get(ctx.UID)

	data, ok := state.Data.(NewTransactionData)
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

	if state.Caller == walletSellCallback {
		amount = -amount
	}

	ctx.ResponseString("Input coin price")

	data.Amount = amount
	ctx.FSM.Update(ctx.UID, State{
		Caller: state.Caller,
		Step:   inputPriceStep,
		Data:   data,
	})
}

func (h *Handler) inputTransactionPriceStep(ctx *Context) {
	state := ctx.FSM.Get(ctx.UID)

	data, ok := state.Data.(NewTransactionData)
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

func (h *Handler) WalletTransactions(ctx *Context) {
	wid, err := strconv.ParseInt(ctx.CallbackDataValue, 10, 64)
	if err != nil {
		ctx.ResponseString("Invalid wallet, expected number")
		return
	}

	wallet, err := h.service.Wallet.Get(ctx.UID, wid)
	if err != nil {
		ctx.ResponseString("Can't find wallet, " + err.Error())
		return
	}

	transactions, err := h.service.Wallet.GetTransactions(wid)
	if err != nil {
		ctx.ResponseString(h.walletError(err))
		return
	}

	if len(transactions) == 0 {
		ctx.ResponseString("No transactions yet")
		return
	}

	spent, earned, profit, err := h.service.Wallet.GetProfit(transactions)
	if err != nil {
		ctx.ResponseString(err.Error())
		return
	}

	ctx.ResponseString(h.formatTransactions(wallet, transactions, spent, earned, profit))
}

func (h *Handler) formatTransactions(wallet models.Wallet, trs []models.Transaction, spent, earned, profit float64) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`
Total spent: %f $
Total earned: %f $
Fixed Profit: %f $
Current Profit: %f $
`, spent, earned, profit, wallet.Balance-math.Abs(profit)))

	for _, tr := range trs {
		s := fmt.Sprintf(`
%s %s:
	Amount: %f
	Price: %f
`, tr.CreatedAt.Format(time.DateTime), tr.Type, tr.Amount, tr.Price)

		if tr.Type == service.WalletBoughtTransaction {
			s += fmt.Sprintf("\tSpend: %f\n", tr.Amount*tr.Price)
		} else if tr.Type == service.WalletSoldTransaction {
			s += fmt.Sprintf("\tEarned: %f\n", -tr.Amount*tr.Price)
		}

		b.WriteString(s)
	}

	return b.String()
}
