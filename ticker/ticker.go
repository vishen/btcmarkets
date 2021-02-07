package ticker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/vishen/btcmarkets/api"
)

func Run(apiKey, privateKey string) error {
	b := api.New(apiKey, privateKey, http.DefaultClient)
	ctx := context.Background()

	balances, err := b.Balance(ctx)
	if err != nil {
		return fmt.Errorf("unable to get balance: %w", err)
	}

	marketBalance := map[string]string{}
	marketIDs := make([]string, 0, len(balances))
	for _, balance := range balances {
		if balance.AssetName == "AUD" {
			continue
		}
		market := balance.AssetName + "-AUD"
		marketIDs = append(marketIDs, market)

		marketBalance[market] = balance.Balance
	}

	marketTickers, err := b.MarketTickers(ctx, marketIDs...)
	if err != nil {
		return fmt.Errorf("unable to get market tickers: %w", err)
	}

	ms := []marketState{}
	for _, m := range marketTickers {
		ms = append(ms, marketState{
			market:    m.MarketID,
			balance:   marketBalance[m.MarketID],
			lastPrice: m.LastPrice,
		})
	}

	sort.Slice(ms, func(i, j int) bool {
		if ms[i].value() > ms[j].value() {
			return true
		}
		//return ms[i].market < ms[j].market
		return false
	})

	for _, m := range ms {
		fmt.Printf("> %s, %s ($%s): $%0.4f\n", m.market, m.balance, m.lastPrice, m.value())
	}

	return nil
}

var state map[string]marketState

type marketState struct {
	market    string
	balance   string
	lastPrice string

	cur  marketData
	prev marketData
}

func (m marketState) value() float64 {
	b, _ := strconv.ParseFloat(m.balance, 64)
	l, _ := strconv.ParseFloat(m.lastPrice, 64)
	return b * l
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
