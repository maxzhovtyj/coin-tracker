package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"testing"
)

func Test(t *testing.T) {
	do, err := binance.NewClient("", "").NewListPricesService().Symbol("BTCUSDT").Do(context.Background(), []binance.RequestOption{}...)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(do[0].Price)
}

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
