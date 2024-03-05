package telegram

type Command string

const (
	startCommand Command = "start"
	// TODO add help command

	// TODO combine newWallet with wallets
	walletsCommand   Command = "wallets"
	newWalletCommand Command = "newWallet"

	buyCommand      Command = "buy"
	sellCommand     Command = "sell"
	netWorthCommand Command = "netWorth"

	// TODO combine subscribe and cancelSubscription into subscriptions
	subscribeCoinCommand          Command = "subscribeCoin"
	cancelCoinSubscriptionCommand Command = "cancelSubscription"

	cancelCommand Command = "cancel"
)

var usefulCommands = []Command{
	walletsCommand,
	newWalletCommand,
	buyCommand,
	sellCommand,
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
	case buyCommand:
		h.NewTransaction(ctx)
	case sellCommand:
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
