package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/config"
	"github.com/maxzhovtyj/coin-tracker/internal/service"
	"go.uber.org/zap"
)

const (
	startMessage     = "start"
	newWalletMessage = "newWallet"
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

func (h *Handler) Init() error {
	h.logger.Infof("authorized on account %s", h.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case startMessage:
			h.CreateUser(&update)
		case newWalletMessage:
			h.NewWallet(&update)
		default:
			h.Response(update.SentFrom().ID, h.UnknownCommand())
		}
	}

	return nil
}
