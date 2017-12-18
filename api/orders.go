package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type OrderHistoryResponse struct {
	Success      bool   `json:"success"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Orders       []struct {
		CreationTime      int64  `json:"creationTime"`
		Currency          string `json:"currency"`
		Instrument        string `json:"instrument"`
		UnformattedPrice  int64  `json:"price"`
		Price             float64
		UnformattedVolume int64 `json:"volume"`
		Volume            float64
		Status            string `json:"status"`
		//Trades       []struct{} `json:"trades"`
	} `json:"orders"`
}

func (c BTCMarketClient) OrderHistory(currency, instrument string) (OrderHistoryResponse, error) {
	orderHistoryRequest := `{"currency":"%s","instrument":"%s","limit":%d,"since":1}`
	data := fmt.Sprintf(orderHistoryRequest, currency, instrument, 100)
	response, err := c.makeRequest("POST", "/order/history", data)
	if err != nil {
		return OrderHistoryResponse{}, err
	}

	orderHistory := OrderHistoryResponse{}
	if err := json.Unmarshal(response, &orderHistory); err != nil {
		return OrderHistoryResponse{}, errors.Wrap(err, "cannot unmarhal order history response")
	}

	for i, o := range orderHistory.Orders {
		orderHistory.Orders[i].Price = FormatInt(o.UnformattedPrice)
		orderHistory.Orders[i].Volume = FormatInt(o.UnformattedVolume)
	}

	return orderHistory, nil

}
