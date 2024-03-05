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

type Context struct {
	UID               int64
	CallbackDataValue string
	Type              messageType
	CommandType       string
	FSM               *FSM
	Update            tgbotapi.Update
	api               *tgbotapi.BotAPI
	logger            *zap.SugaredLogger
}

func NewContext(update tgbotapi.Update, api *tgbotapi.BotAPI, fsm *FSM, logger *zap.SugaredLogger) *Context {
	var cmdType string

	msgType := RegularMessage

	if update.Message != nil && update.Message.IsCommand() {
		msgType = CommandMessage
		cmdType = update.Message.Command()
	} else if update.CallbackQuery != nil {
		msgType = CallbackMessage
	}

	return &Context{
		UID:         update.SentFrom().ID,
		Update:      update,
		Type:        msgType,
		CommandType: cmdType,
		FSM:         fsm,
		api:         api,
		logger:      logger,
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

func (ctx *Context) ResponseWithError(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	return ctx.api.Send(msg)
}

func (ctx *Context) UnknownCommand() string {
	return "Sorry, I don't what command it is"
}
