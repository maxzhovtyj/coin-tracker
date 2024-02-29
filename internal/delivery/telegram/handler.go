package telegram

import (
	"fmt"
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
}

func NewHandler(cfg *config.Config, service *service.Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		cfg:     cfg,
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Init() error {
	bot, err := tgbotapi.NewBotAPI(h.cfg.TelegramApiToken)
	if err != nil {
		return fmt.Errorf("failed to get telegram api token: %w", err)
	}

	h.logger.Infof("authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		var msg tgbotapi.MessageConfig

		switch update.Message.Command() {
		case startMessage:
			msg = h.CreateUser(&update)
		case newWalletMessage:
			msg = h.NewWallet(&update)
		default:
		}

		if _, err = bot.Send(msg); err != nil {
			h.logger.Errorf("failed to sent response message: %w", err)
		}
	}

	return nil
}
