package telegram

import (
	"strings"
)

const (
	walletCallback             = "wallet"
	newTransactionCallback     = "newTransaction"
	walletNewCallback          = "walletNew"
	walletDeleteCallback       = "walletDelete"
	walletTransactionsCallback = "walletTransactions"
	walletBuyCallback          = "walletBuy"
	walletSellCallback         = "walletSell"
)

func (h *Handler) Callbacks(ctx *Context) {
	cbData := ctx.Update.CallbackData()

	idx := strings.IndexRune(cbData, '=')
	if idx == -1 || len(cbData) == idx {
		h.logger.Errorf("invalid callback data '%s'", ctx.Update.CallbackData())
		return
	}

	cbName := cbData[:idx]
	ctx.CallbackName = cbName
	ctx.CallbackDataValue = cbData[idx+1:]

	switch cbName {
	case walletCallback:
		h.Wallet(ctx)
	case walletNewCallback:
		h.NewWallet(ctx)
	case walletDeleteCallback:
		h.DeleteWallet(ctx)
	case walletBuyCallback, walletSellCallback:
		h.NewTransaction(ctx)
	case newTransactionCallback:
		h.ResolveNewTransactionSteps(ctx)
	case walletTransactionsCallback:
		h.WalletTransactions(ctx)
	}

}
