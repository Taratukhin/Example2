package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	binance "github.com/aiviaio/go-binance/v2"
)

var Client *binance.Client

func GetPrice(symbol string, ch chan map[int]binance.SymbolPrice) {
	result, err := Client.NewListPricesService().Symbol(symbol).Do(context.Background())

	if err != nil {
		fmt.Printf("GetPrice: error to get %v\n", err)
		return
	}

	answer := make(map[int]binance.SymbolPrice)

	answer[0] = *result[0]

	ch <- answer
}

func main() {
	Client = binance.NewClient("", "")

	wg := sync.WaitGroup{}

	symbols, err := Client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	result := make([]string, 0, 5)

	for i := 0; i < 5 && i < len(symbols.Symbols); i++ {
		result = append(result, symbols.Symbols[i].Symbol)
	}

	// I don’t know why to transfer the card through the channel, but in the conditions of this task it says so
	ch := make(chan map[int]binance.SymbolPrice)

	for _, symbol := range result {
		wg.Add(1)
		go GetPrice(symbol, ch)
	}

	go func() {
		for p := range ch {
			fmt.Println(p[0].Symbol, p[0].Price)
			wg.Done()
		}
	}()

	wg.Wait()
}
