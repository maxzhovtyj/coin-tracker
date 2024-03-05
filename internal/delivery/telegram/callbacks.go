package telegram

import (
	"strings"
)

const (
	walletCallback             = "wallet"
	newTransactionCallback     = "newTransaction"
	walletTransactionsCallback = "walletTransactions"
)

func (h *Handler) Callbacks(ctx *Context) {
	cbData := ctx.Update.CallbackData()

	idx := strings.IndexRune(cbData, '=')
	if idx == -1 || len(cbData) == idx {
		h.logger.Errorf("invalid callback data '%s'", ctx.Update.CallbackData())
		return
	}

	cbName := cbData[:idx]
	ctx.CallbackDataValue = cbData[idx+1:]

	switch cbName {
	case walletCallback:
		h.Wallet(ctx)
	case newTransactionCallback:
		h.ResolveNewTransactionSteps(ctx)
	case walletTransactionsCallback:
		h.WalletTransactions(ctx)
	}

}
