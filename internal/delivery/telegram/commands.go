package telegram

type Command string

const (
	startCommand                  Command = "start"
	walletsCommand                Command = "wallets"
	newWalletCommand              Command = "newWallet"
	newTransactionCommand         Command = "buy"
	netWorthCommand               Command = "netWorth"
	subscribeCoinCommand          Command = "subscribeCoin"
	cancelCoinSubscriptionCommand Command = "cancelSubscription"
	cancelCommand                 Command = "cancel"
)

var usefulCommands = []Command{
	walletsCommand,
	newWalletCommand,
	newTransactionCommand,
	netWorthCommand,
	subscribeCoinCommand,
	cancelCoinSubscriptionCommand,
	cancelCommand,
}

func (h *Handler) Commands(ctx *Context) {
	switch Command(ctx.Command()) {
	case startCommand:
		h.CreateUser(ctx)
	case walletsCommand:
		h.Wallets(ctx)
	case newWalletCommand:
		h.NewWallet(ctx)
	case newTransactionCommand:
		h.NewTransaction(ctx)
	case netWorthCommand:
		h.UserNetWorth(ctx)
	case subscribeCoinCommand:
		h.SubscribeCoin(ctx)
	case cancelCoinSubscriptionCommand:
		h.CancelCoinSubscription(ctx)
	case cancelCommand:
		h.Cancel(ctx)
	default:
		ctx.ResponseString(ctx.UnknownCommand())
	}
}

func (h *Handler) Cancel(ctx *Context) {
	ctx.FSM.Remove(ctx.UID)
	ctx.ResponseString("All of your previous commands are canceled")
}
