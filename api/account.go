package api

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type AccountBalance []struct {
	BalanceUnformatted int64 `json:"balance"`
	Balance            float64
	PendingFunds       int64  `json:"pendingFunds"`
	Currency           string `json:"currency"`
}

func (c BTCMarketClient) AccountBalance() (AccountBalance, error) {
	response, err := c.makeRequest("GET", "/account/balance", "")
	if err != nil {
		return AccountBalance{}, err
	}

	accountBalance := AccountBalance{}
	if err := json.Unmarshal(response, &accountBalance); err != nil {
		return AccountBalance{}, errors.Wrap(err, "cannot unmarshal account balance")
	}

	for i, ab := range accountBalance {
		accountBalance[i].Balance = FormatInt(ab.BalanceUnformatted)
	}

	return accountBalance, nil
}
