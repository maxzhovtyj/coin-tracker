package models

type Wallet struct {
	Id      int64   `json:"id"`
	UserID  int64   `json:"user_id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price,omitempty"`
	Amount  float64 `json:"amount,omitempty"`
	Balance float64 `json:"balance,omitempty"`
}
