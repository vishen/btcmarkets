package api

import "encoding/base64"

const (
	baseUrl = "https://api.btcmarkets.net"
)

type BTCMarketClient struct {
	publicAPIKey  string
	privateAPIKey []byte
}

func NewBTCMarketClient(publicAPIKey, privateAPIKeyB64Encoded string) (*BTCMarketClient, error) {

	client := BTCMarketClient{
		publicAPIKey: publicAPIKey,
	}
	var err error

	if client.privateAPIKey, err = base64.StdEncoding.DecodeString(privateAPIKeyB64Encoded); err != nil {
		return nil, err
	}

	return &client, nil

}

func (c BTCMarketClient) makeRequest(requestType, path, data string) ([]byte, error) {
	return CustomRequest(requestType, path, data, c.publicAPIKey, c.privateAPIKey)
}
