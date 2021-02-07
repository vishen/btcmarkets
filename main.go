package main

import (
	"flag"
	"fmt"

	"github.com/vishen/btcmarkets/ticker"
)

var (
	apiKeyFlag     = flag.String("api-key", "", "BTC Markets API Key")
	privateKeyFlag = flag.String("private-key", "", "BTC Markets Private Key base64 encoded")
)

func main() {
	flag.Parse()

	switch "" {
	case *apiKeyFlag, *privateKeyFlag:
		fmt.Printf("requires -api-key and -private-key\n")
		return
	}

	ticker.Run(*apiKeyFlag, *privateKeyFlag)
}
