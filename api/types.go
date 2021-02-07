package api

type BalanceResponse []struct {
	AssetName string `json:"assetName"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Locked    string `json:"locked"`
}

type MarketTickersResponse []MarketTicker

type MarketTicker struct {
	MarketID  string `json:"marketId"`
	BestBid   string `json:"bestBid"`
	BestAsk   string `json:"bestAsk"`
	LastPrice string `json:"lastPrice"`
	Volume24h string `json:"volume24h"`
	Price24h  string `json:"price24h"`
	Low24h    string `json:"low24h"`
	High24h   string `json:"high24h"`
	Timestamp string `json:"timestamp"`
}

type MarketCandleResponse []MarketCandle

type MarketCandle struct {
	Timestamp string
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
}
