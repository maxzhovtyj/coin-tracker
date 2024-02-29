package app

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/maxzhovtyj/coin-tracker/internal/config"
	"github.com/maxzhovtyj/coin-tracker/internal/delivery/telegram"
	"github.com/maxzhovtyj/coin-tracker/internal/service"
	"github.com/maxzhovtyj/coin-tracker/internal/storage"
	"github.com/maxzhovtyj/coin-tracker/pkg/log/applogger"
)

func Run() error {
	logger := applogger.New()

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("can't init config: %w", err)
	}

	logger.Infof("start application")

	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{})

	db, err := sql.Open("sqlite3_custom", cfg.DBPath)
	if err != nil {
		logger.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	st := storage.New(db)
	s := service.New(st)
	handler := telegram.NewHandler(cfg, s, logger)
	err = handler.Init()
	if err != nil {
		return fmt.Errorf("failed to init telegram handler: %w", err)
	}

	return nil
}
