package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	apiURL = "https://api.btcmarkets.net"
)

/*
Rate Limits:

The default rate limit (unless specified) for public and private APIs is up to 50 calls per 10 seconds.

Exceptions:

- order placement API	/v3/orders POST	30 calls per 10 seconds
- batch order API	/v3/batchorders POST	5 calls per 10 seconds
- withdraw request API	/v3/withdrawals POST	10 calls per 10 seconds
- creating new report API	/v3/reports POST	1 call per 10 seconds
*/

type BTCMarketAPI struct {
	apiKey     string
	privateKey string

	httpClient *http.Client
}

// privateKey should be base64 encoded
func New(apiKey, privateKey string, httpClient *http.Client) BTCMarketAPI {
	return BTCMarketAPI{
		apiKey:     apiKey,
		privateKey: privateKey,
		httpClient: httpClient,
	}
}

type CandleWindow string

const (
	CandleWindow_1m CandleWindow = "1m"
	CandleWindow_1h CandleWindow = "1h"
	CandleWindow_1d CandleWindow = "1d"
)

func (b BTCMarketAPI) MarketCandles(ctx context.Context, marketID string, from, to time.Time, window CandleWindow) (MarketCandleResponse, error) {
	/*

		[time,open,high,low,close,volume]
		[
			[
			"2019-09-02T18:00:00.000000Z",
			"15100",
			"15200",
			"15100",
			"15199",
			"4.11970335"
			],
			[
			"2019-09-02T17:00:00.000000Z",
			"14879.75",
			"15115",
			"14861.99",
			"15115",
			"10.01840031"
			]
		]
	*/

	// Switch if from and to are around the wrong way
	if to.Before(from) {
		to, from = from, to
	}

	path := "/v3/markets/" + marketID + "/candles?timeWindow=" + string(window) + "&from=" + ISO8601(from) + "&to=" + ISO8601(to)
	fmt.Println("PATH: ", path)

	var r MarketCandleResponse
	data, err := b.httpRequest(ctx, "GET", path, "")
	if err != nil {
		return r, err
	}

	resp := [][6]string{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return r, fmt.Errorf("error decoding response: %w", err)
	}

	for _, a := range resp {
		r = append(r, MarketCandle{
			Timestamp: a[0],
			Open:      a[1],
			High:      a[2],
			Low:       a[3],
			Close:     a[4],
			Volume:    a[5],
		})
	}

	return r, nil
}

func (b BTCMarketAPI) MarketTickers(ctx context.Context, marketIDs ...string) (MarketTickersResponse, error) {
	path := "/v3/markets/tickers?"

	for i, mid := range marketIDs {
		path += "marketId=" + mid
		if i < len(marketIDs)-1 {
			path += "&"
		}
	}

	var r MarketTickersResponse
	data, err := b.httpRequest(ctx, "GET", path, "")
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return r, fmt.Errorf("error decoding response: %w", err)
	}
	return r, nil
}

func (b BTCMarketAPI) Balance(ctx context.Context) (BalanceResponse, error) {
	path := "/v3/accounts/me/balances"

	var r BalanceResponse
	data, err := b.httpRequest(ctx, "GET", path, "")
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return r, fmt.Errorf("error decoding response: %w", err)
	}
	return r, nil
}

func (b BTCMarketAPI) httpRequest(ctx context.Context, method, path, body string) ([]byte, error) {
	url, err := url.Parse(apiURL + path)
	if err != nil {
		return nil, fmt.Errorf("error parsing url %q: %w", apiURL+path, err)
	}
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, _ := http.NewRequestWithContext(ctx, method, url.String(), bodyReader)
	req.Header = b.buildAuthHeaders("GET", url.Path, body)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling btc market api: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// fmt.Printf("Response Body: %s\n", responseBody)

	// If there was an error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("headers: ", req.Header)
		fmt.Println("url: ", url)
		return nil, fmt.Errorf("error with request: status_code=%d response=%q", resp.StatusCode, responseBody)
	}
	return responseBody, nil
}

func (b BTCMarketAPI) buildAuthHeaders(method string, path string, body string) http.Header {
	nowMs := ts(time.Now())
	stringToSign := method + path + nowMs
	if body != "" {
		stringToSign += body
	}
	return http.Header{
		"Content-Type":      []string{"application/json"},
		"Accept":            []string{"application/json"},
		"Accept-Charset":    []string{"UTF-8"},
		"BM-AUTH-APIKEY":    []string{b.apiKey},
		"BM-AUTH-TIMESTAMP": []string{nowMs},
		"BM-AUTH-SIGNATURE": []string{signMessage(b.privateKey, stringToSign)},
	}
}

func signMessage(key, message string) string {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Fatal(err)
	}
	mac := hmac.New(sha512.New, decodedKey)
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func ts(t time.Time) string {
	return strconv.FormatInt(t.UTC().UnixNano()/1000000, 10)
}

func ISO8601(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
