package main

import (
	"github.com/maxzhovtyj/coin-tracker/internal/app"
	logger "github.com/maxzhovtyj/coin-tracker/pkg/log/applogger"
)

func main() {
	log := logger.New()
	err := app.Run()
	if err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
