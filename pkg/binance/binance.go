package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"time"
)

type API interface {
	Info(symbol string) (*binance.SymbolPrice, error)
}

type api struct {
	client *binance.Client
}

func NewAPI() API {
	return &api{
		client: binance.NewClient("", ""),
	}
}

func (a *api) Info(symbol string) (*binance.SymbolPrice, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	symbols, err := a.client.NewListPricesService().Symbol(symbol).Do(ctx, []binance.RequestOption{}...)
	if err != nil {
		return nil, err
	}

	if len(symbols) == 0 {
		return nil, fmt.Errorf("empty result")
	}

	return symbols[0], nil
}
