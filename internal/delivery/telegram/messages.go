package telegram

func (h *Handler) Messages(ctx *Context) {
	state := ctx.FSM.Get(ctx.UID)

	switch state.Caller {
	case walletNewCallback:
		h.ResolveNewWalletSteps(ctx)
	case walletBuyCallback, walletSellCallback:
		h.ResolveNewTransactionSteps(ctx)
	case subscribeCoinCommand.String():
		h.ResolveSubscribeCoinSteps(ctx)
	}
}
