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
