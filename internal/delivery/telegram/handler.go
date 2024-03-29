package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/config"
	"github.com/maxzhovtyj/coin-tracker/internal/service"
	"go.uber.org/zap"
	"strings"
)

type Handler struct {
	cfg     *config.Config
	service *service.Service
	logger  *zap.SugaredLogger
	bot     *tgbotapi.BotAPI
}

func NewHandler(cfg *config.Config, bot *tgbotapi.BotAPI, service *service.Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		cfg:     cfg,
		service: service,
		logger:  logger,
		bot:     bot,
	}
}

var usersWhitelist = map[string]struct{}{
	"maxzhovtyj":  {},
	"vadimkrivko": {},
	"Andrii2324":  {},
}

func (h *Handler) Init() error {
	h.logger.Infof("authorized on account %s", h.bot.Self.UserName)

	h.logger.Infof("init subscriptions worker")
	go h.SubscriptionsWorker()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	fsm := NewFSM()

	updates := h.bot.GetUpdatesChan(u)
	for update := range updates {
		ctx := NewContext(update, h.bot, fsm, h.logger)

		if _, ok := usersWhitelist[update.SentFrom().UserName]; !ok {
			ctx.ResponseString("Sorry, you are not allowed to use this bot")
			continue
		}

		switch ctx.Type {
		case CallbackMessage:
			h.Callbacks(ctx)
		case CommandMessage:
			h.Commands(ctx)
		case RegularMessage:
			h.Messages(ctx)
		}
	}

	return nil
}

func (h *Handler) resolveWalletName(n string) string {
	n = strings.TrimSpace(n)
	n = strings.ToUpper(n)
	return n
}
