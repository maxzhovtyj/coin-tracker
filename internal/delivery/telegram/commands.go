package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Command string

func (c Command) String() string {
	return string(c)
}

const (
	startCommand         Command = "start"
	walletsCommand       Command = "wallets"
	netWorthCommand      Command = "netWorth"
	subscriptionsCommand Command = "subscriptions"
	cancelCommand        Command = "cancel"
)

var keyboardMarkup = [][]tgbotapi.KeyboardButton{
	{tgbotapi.NewKeyboardButton("/" + walletsCommand.String()), tgbotapi.NewKeyboardButton("/" + netWorthCommand.String())},
	{tgbotapi.NewKeyboardButton("/" + subscriptionsCommand.String())},
	{tgbotapi.NewKeyboardButton("/" + cancelCommand.String())},
}

func (h *Handler) Commands(ctx *Context) {
	switch Command(ctx.Command()) {
	case startCommand:
		h.CreateUser(ctx)
	case walletsCommand:
		h.Wallets(ctx)
	case netWorthCommand:
		h.UserNetWorth(ctx)
	case subscriptionsCommand:
		h.Subscriptions(ctx)
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
