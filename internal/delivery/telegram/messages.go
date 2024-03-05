package telegram

func (h *Handler) Messages(ctx *Context) {
	state := ctx.FSM.Get(ctx.UID)

	switch state.Command {
	case newWalletCommand:
		h.ResolveNewWalletSteps(ctx)
	case buyCommand, sellCommand:
		h.ResolveNewTransactionSteps(ctx)
	case subscribeCoinCommand:
		h.ResolveSubscribeCoinSteps(ctx)
	}
}
