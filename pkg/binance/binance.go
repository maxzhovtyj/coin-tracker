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
	CoinTicker(coin string, windowSize string) (SymbolTicker, error)
}

type api struct {
	client *binance.Client
}

func NewAPI() API {
	return &api{
		client: binance.NewClient("", ""),
	}
}

type Coin struct {
	Symbol string
	Price  float64
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

type SymbolTicker struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	LastPrice          string `json:"lastPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstId            int64  `json:"firstId"`
	LastId             int64  `json:"lastId"`
	Count              int64  `json:"count"`
}

func (a *api) CoinTicker(coin string, windowSize string) (SymbolTicker, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	tickers, err := a.client.NewListSymbolTickerService().Symbol(coin).WindowSize(windowSize).Do(ctx)
	if err != nil {
		return SymbolTicker{}, err
	}

	if len(tickers) == 0 {
		return SymbolTicker{}, fmt.Errorf("empty coins result")
	}

	t := tickers[0]

	return SymbolTicker{
		Symbol:             t.Symbol,
		PriceChange:        t.PriceChange,
		PriceChangePercent: t.PriceChangePercent,
		WeightedAvgPrice:   t.WeightedAvgPrice,
		OpenPrice:          t.OpenPrice,
		HighPrice:          t.HighPrice,
		LowPrice:           t.LowPrice,
		LastPrice:          t.LastPrice,
		Volume:             t.Volume,
		QuoteVolume:        t.QuoteVolume,
		OpenTime:           t.OpenTime,
		CloseTime:          t.CloseTime,
		FirstId:            t.FirstId,
		LastId:             t.LastId,
		Count:              t.Count,
	}, nil
}
