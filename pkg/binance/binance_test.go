package binance

import (
	"fmt"
	"testing"
)

func TestApi_CoinsList(t *testing.T) {
	s := []string{"ETHUSDT", "BTCUSDT"}

	list, err := NewAPI().CoinsList(s...)
	if err != nil {
		t.Error(err)
	}

	for _, l := range list {
		fmt.Println(l.Symbol, l.Price)
	}
}

func TestApi_CoinTicker(t *testing.T) {
	ticker, err := NewAPI().CoinTicker("APTUSDT", "1h")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(ticker)
}
