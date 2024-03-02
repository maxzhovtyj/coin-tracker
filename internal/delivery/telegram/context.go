package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strings"
)

type messageType string

const (
	CallbackMessage messageType = "callback"
	CommandMessage  messageType = "command"
	RegularMessage  messageType = "message"
)

type State struct {
	Command string
}

var usersState = map[int64]*State{}

type Context struct {
	UID               int64
	CallbackDataValue string
	Update            tgbotapi.Update
	Type              messageType
	api               *tgbotapi.BotAPI
	logger            *zap.SugaredLogger
}

func NewContext(update tgbotapi.Update, api *tgbotapi.BotAPI, logger *zap.SugaredLogger) *Context {
	msgType := RegularMessage

	if update.Message != nil && update.Message.IsCommand() {
		msgType = CommandMessage
	} else if update.CallbackQuery != nil {
		msgType = CallbackMessage
	}

	if msgType == RegularMessage {

	}

	return &Context{
		UID:    update.SentFrom().ID,
		Update: update,
		Type:   msgType,
		api:    api,
		logger: logger,
	}
}

func (ctx *Context) CommandArgs() []string {
	return strings.Split(ctx.Update.Message.CommandArguments(), " ")
}

func (ctx *Context) Command() string {
	return ctx.Update.Message.Command()
}

func (ctx *Context) ResponseString(msg string) {
	_, err := ctx.api.Send(tgbotapi.NewMessage(ctx.UID, msg))
	if err != nil {
		ctx.logger.Errorf("failed to send message: %v", err)
		return
	}
}

func (ctx *Context) Response(msg tgbotapi.MessageConfig) {
	_, err := ctx.api.Send(msg)
	if err != nil {
		ctx.logger.Errorf("failed to send message: %v", err)
		return
	}
}

func (ctx *Context) UnknownCommand() string {
	return "Sorry, I don't what command it is"
}
