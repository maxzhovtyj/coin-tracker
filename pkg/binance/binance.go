package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"strconv"
	"time"
)

type API interface {
	Coin(symbol string) (float64, error)
	CoinsList(symbols ...string) ([]Coin, error)
}

type api struct {
	client *binance.Client
}

func NewAPI() API {
	return &api{
		client: binance.NewClient("", ""),
	}
}

func (a *api) Coin(symbol string) (float64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	symbols, err := a.client.NewListPricesService().Symbol(symbol).Do(ctx, []binance.RequestOption{}...)
	if err != nil {
		return 0, err
	}

	if len(symbols) == 0 {
		return 0, fmt.Errorf("empty result")
	}

	price, err := strconv.ParseFloat(symbols[0].Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

type Coin struct {
	Symbol string
	Price  float64
}

func (a *api) CoinsList(symbols ...string) ([]Coin, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	list, err := a.client.NewListPricesService().Symbols(symbols).Do(ctx)
	if err != nil {
		return nil, err
	}

	coins := make([]Coin, len(list))
	for i, s := range list {
		coins[i].Symbol = s.Symbol
		coins[i].Price, err = strconv.ParseFloat(s.Price, 64)
		if err != nil {
			return nil, err
		}
	}

	return coins, nil
}
