package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func GetTimestamp() int64 {
	return time.Now().Unix() * 1000
}

func FormatInt(value int64) float64 {
	return float64(value) / float64(math.Pow(10, 8))
}

func CustomRequest(requestType, path, data, publicKey string, privateKey []byte) ([]byte, error) {

	// Get the current timestamp as a string
	timestamp := strconv.FormatInt(GetTimestamp(), 10)

	// Sign this string which contains the info about the request
	stringToSign := path + "\n" + timestamp + "\n" + data
	stringToSign = fmt.Sprintf("%s\n%s\n%s", path, timestamp, data)

	// Sign the 'stringToSign' using sha512 and the private api key
	mac := hmac.New(sha512.New, privateKey)
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// Create a request object
	req, err := http.NewRequest(requestType, baseUrl+path, bytes.NewBufferString(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to build NewRequest")
	}

	// Add the headers to the request object
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Golang BTCMarkets client")
	req.Header.Set("accept-charset", "utf-8")
	req.Header.Set("apikey", publicKey)
	req.Header.Set("signature", signature)
	req.Header.Set("timestamp", timestamp)

	// Make the request
	// TODO(vishen): Allow custom http client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to make request")
	}

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't read response body: status code %d", resp.StatusCode)
	}
	resp.Body.Close()

	return responseBody, nil
}
