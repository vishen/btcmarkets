package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type MarketPrice struct {
	BestBid    float64 `json:"bestBid"`
	BestAsk    float32 `json:"bestAsk"`
	LastPrice  float64 `json:"lastPrice"`
	Currency   string  `json:"currency"`
	Instrument string  `json:"instrument"`
	Timestamp  int64   `json:"timestamp"`
	Volume24h  float32 `json:"volume24h"`
}

func (c BTCMarketClient) MarketTick(t1, t2 string) (MarketPrice, error) {
	url := fmt.Sprintf("/market/%s/%s/tick", t1, t2)
	response, err := c.makeRequest("GET", url, "")
	if err != nil {
		return MarketPrice{}, err
	}

	marketPrice := MarketPrice{}
	if err := json.Unmarshal(response, &marketPrice); err != nil {
		return MarketPrice{}, errors.Wrap(err, "cannot unmarshal market price")
	}

	return marketPrice, nil

}
