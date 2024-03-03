package models

import (
	"time"
)

type CoinSubscriptionData struct {
	Interval string
	CoinName string
}

type Subscription struct {
	ID             int64
	Type           string
	ChatID         int64
	Data           string
	NotifyInterval time.Duration
	LastNotifiedAt time.Time
}
