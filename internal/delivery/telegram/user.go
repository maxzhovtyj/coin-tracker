package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	"strings"
)

func (h *Handler) CreateUser(ctx *Context) {
	msgText := "Congratulations! New user successfully created"

	_, err := h.service.User.Create(ctx.UID)
	if errors.Is(err, models.ErrUserAlreadyExists) {
		msgText = "Sorry, such user already exists"
	} else if err != nil {
		ctx.ResponseString(h.startError(err))
		return
	}

	msg := tgbotapi.NewMessage(ctx.UID, msgText)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(keyboardMarkup...)

	ctx.Response(msg)
}

func (h *Handler) startError(err error) string {
	return fmt.Sprintf("Sorry, I can't create new user, reason: %v", err)
}

func (h *Handler) UserNetWorth(ctx *Context) {
	netWorth, err := h.service.NetWorth(ctx.UID)
	if err != nil {
		ctx.ResponseString("Sorry, can't get your net worth, " + err.Error())
		return

	}

	ctx.ResponseString(h.netWorthSuccess(netWorth))
}

func (h *Handler) netWorthSuccess(worth models.UserNetWorth) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Your net worth: %f $ \n\n", worth.Balance))
	for _, w := range worth.Wallets {
		b.WriteString(fmt.Sprintf("* Name: %s\n  Price: %f\n  Amount: %f\n  Balance: %f\n\n", w.Name, w.Price, w.Amount, w.Balance))
	}

	return b.String()
}
