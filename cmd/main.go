package main

import (
	"github.com/maxzhovtyj/coin-tracker/internal/app"
	"github.com/maxzhovtyj/coin-tracker/pkg/logger"
)

func main() {
	err := app.Run()
	if err != nil {
		logger.Fatalf("can't run application: %v", err)
	}
}
