package ticker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vishen/btcmarkets/api"
)

func Run(apiKey, privateKey string) {
	b := api.New(apiKey, privateKey, http.DefaultClient)
	ctx := context.Background()

	balances, err := b.Balance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	marketIDs := make([]string, 0, len(balances))
	for _, balance := range balances {
		fmt.Printf("> %s (%s) | available=%s locked=%s\n", balance.AssetName, balance.Balance, balance.Available, balance.Locked)
		if balance.AssetName != "AUD" {
			market := balance.AssetName + "-AUD"
			marketIDs = append(marketIDs, market)
		}

	}
	_ = marketIDs

	/*
		marketTickerT := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-marketTickerT.C:
				marketTickers, err := b.MarketTickers(ctx, marketIDs...)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	*/
}

var state map[string]marketState

type marketState struct {
	market  string
	balance string

	cur  marketData
	prev marketData
}

type marketData struct {
	bestBid   string
	bestAsk   string
	lastPrice string

	high5m   string
	low5m    string
	volume5m string

	high10m   string
	low10m    string
	volume10m string
}

func marketCandles(ctx context.Context, b api.BTCMarketAPI) {
	now := time.Now()

	/*
		now5m := now.Add(time.Minute * -5)
		now1h := now.Add(time.Minute * -60)
		marketCandles, err := b.MarketCandles(ctx, "ETH-AUD", now1h, now5m, api.CandleWindow_1m)
		if err != nil {
			log.Fatal(err)
		}

		max := len(marketCandles)
		candles := map[string]api.MarketCandle{
			"60m": marketCandles[0],
			"30m": marketCandles[max/2-max/8],
			"10m": marketCandles[max/2+max/4],
			"5m":  marketCandles[max-1],
		}

		fmt.Println("60m", candles["60m"])
		fmt.Println("30m", candles["30m"])
		fmt.Println("10m", candles["10m"])
		fmt.Println("5m", candles["5m"])
	*/

	now6h := now.Add(time.Hour * -6)
	now48h := now.Add(time.Hour * -48)
	marketCandles, err := b.MarketCandles(ctx, "ETH-AUD", now6h, now48h, api.CandleWindow_1h)
	if err != nil {
		log.Fatal(err)
	}

	for i, mc := range marketCandles {
		fmt.Printf("%d) %#v\n", i, mc)
	}

	/*
		max := len(marketCandles)
		candles := map[string]api.MarketCandle{
			"48h": marketCandles[0],
			"24h": marketCandles[max/2-max/8],
			"12h": marketCandles[max/2],
			"6h":  marketCandles[max-1],
		}

		fmt.Println("6h", candles["6h"])
		fmt.Println("12h", candles["12h"])
		fmt.Println("24h", candles["24h"])
		fmt.Println("48h", candles["48h"])
	*/
	return

}
