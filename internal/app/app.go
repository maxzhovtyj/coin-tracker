package app

import (
	"fmt"
	"github.com/maxzhovtyj/coin-tracker/internal/config"
	"github.com/maxzhovtyj/coin-tracker/internal/delivery/telegram"
	"github.com/maxzhovtyj/coin-tracker/pkg/log/applogger"
)

func Run() error {
	logger := applogger.New()

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("can't init config: %w", err)
	}

	logger.Infof("start application")

	handler := telegram.NewHandler(cfg, logger)
	err = handler.Init()
	if err != nil {
		return fmt.Errorf("failed to init telegram handler: %w", err)
	}

	return nil
}
