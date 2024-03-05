package models

import "time"

type Transaction struct {
	ID        int64
	WalletID  int64
	Amount    float64
	Price     float64
	CreatedAt time.Time
}
