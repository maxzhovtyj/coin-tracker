package main

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"log"
)

func main() {
	cl := binance.NewClient("", "")
	do, err := cl.NewListPricesService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, sp := range do {
		fmt.Println(sp.Symbol, sp.Price)
	}
}
