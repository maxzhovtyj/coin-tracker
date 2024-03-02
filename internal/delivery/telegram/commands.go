package telegram

type Command string

const (
	startCommand          Command = "start"
	walletsCommand        Command = "wallets"
	newWalletCommand      Command = "newWallet"
	newTransactionCommand Command = "buy"
)

var usefulCommands = []Command{
	walletsCommand,
	newWalletCommand,
	newTransactionCommand,
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
	default:
		ctx.ResponseString(ctx.UnknownCommand())
	}
}
