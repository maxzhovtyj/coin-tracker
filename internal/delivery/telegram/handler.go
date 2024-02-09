package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/config"
	"go.uber.org/zap"
)

const (
	startMessage = "start"
	getWallets   = "wallets"
)

type Handler struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func NewHandler(cfg *config.Config, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: logger,
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

		user := update.SentFrom()

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case startMessage:
			msg.Text = fmt.Sprintf("Hello, %s, thank you for using me", user.FirstName)
		case getWallets:
			msg.Text = "Your wallets:"
		default:
			msg.Text = "I don't know that command"
		}

		if _, err = bot.Send(msg); err != nil {
			h.logger.Errorf("failed to sent response message: %w", err)
		}
	}

	return nil
}
