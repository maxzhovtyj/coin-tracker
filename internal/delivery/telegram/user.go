package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	"math"
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
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(getCommandsKeyboard(usefulCommands)...)

	ctx.Response(msg)
}

func (h *Handler) startError(err error) string {
	return fmt.Sprintf("Sorry, I can't create new user, reason: %v", err)
}

// TODO Refactor keyboards

func getCommandsKeyboard(commands []Command) [][]tgbotapi.KeyboardButton {
	rows := math.Ceil(float64(len(commands)) / float64(2))
	keyboard := make([][]tgbotapi.KeyboardButton, int(rows))

	var row int
	var col int

	for _, t := range commands {
		keyboard[row] = append(keyboard[row], tgbotapi.NewKeyboardButton("/"+string(t)))

		if (col+1)%2 == 0 {
			col = 0
			row++
		} else {
			col++
		}
	}

	return keyboard
}
