package main

import (
	"fmt"
	"log"
	"os"

	"github.com/micro/cli"
	"github.com/vishen/btcmarkets/api"
)

const (
	version      = "0.0.1"
	mainCurrency = "AUD"
)

var (
	publicAPIKey  string
	privateAPIKey string
)

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Name = "BTC Markets"
	app.HelpName = "BTC Markets"
	app.Usage = "Account and market prices for BTC Markets cryptocurrencies"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "public_key",
			Usage:  "Public api key",
			EnvVar: "BTC_MARKETS_PUBLIC_API_KEY",
		},
		cli.StringFlag{
			Name:   "private_key",
			Usage:  "Base64 encoded private key",
			EnvVar: "BTC_MARKETS_PRIVATE_API_KEY",
		},
	}
	app.Before = func(ctx *cli.Context) error {
		publicAPIKey = ctx.GlobalString("public_key")
		privateAPIKey = ctx.GlobalString("private_key")
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "account",
			Usage: "Account balance and analysis",
			Action: func(c *cli.Context) {
				accountAnalysis()
			},
		},
	}
	app.Run(os.Args)
}

func accountAnalysis() {

	client, err := api.NewBTCMarketClient(publicAPIKey, privateAPIKey)
	if err != nil {
		log.Fatal(err)
	}

	accountBalance, err := client.AccountBalance()
	if err != nil {
		log.Fatal(err)
	}

	gtotal := 0.0
	gActualTotal := 0.0
	for _, b := range accountBalance {
		fmt.Printf("############ %s #############\n", b.Currency)
		fmt.Printf("Balance %0.3f%s\n", b.Balance, b.Currency)
		if b.Currency == mainCurrency {
			fmt.Println()
			continue
		}

		marketPrice, err := client.MarketTick(b.Currency, mainCurrency)
		if err != nil {
			log.Printf("Error getting market price for '%s': %s\n", b.Currency, err)
			continue
		}

		orderHistory, err := client.OrderHistory(mainCurrency, b.Currency)
		if err != nil {
			log.Printf("Error getting order history for '%s': %s\n", b.Currency, err)
			continue
		}

		fmt.Printf(
			"Current market price: last_price=%0.3f, best_bid=%0.3f, best_ask=%0.3f\n",
			marketPrice.LastPrice, marketPrice.BestBid, marketPrice.BestAsk,
		)

		ototal := 0.0
		oActualTotal := 0.0
		for _, o := range orderHistory.Orders {
			if o.Status == "Fully Matched" {
				if ototal == 0 {
					fmt.Printf("Transactions:\n")
				}
				price := o.Price
				volume := o.Volume
				total := price * volume
				actualTotal := marketPrice.LastPrice * volume
				ototal += total
				oActualTotal += actualTotal
				var status string
				if total <= actualTotal {
					status = "PROFIT"
				} else {
					status = "LOSS"
				}

				fmt.Printf(
					">> %0.3f%s * %0.3f = %0.3f%s \t| current value: %0.3f%s \t| %s %0.3f\n",
					price, o.Currency, volume, total, o.Currency,
					actualTotal, o.Currency,
					status, actualTotal-total,
				)
			}
		}

		var status string
		if ototal <= oActualTotal {
			status = "PROFIT"
		} else {
			status = "LOSS"
		}
		if ototal > 0 {
			fmt.Printf(
				"%s -> total_spend=%0.3f, current_worth=%0.3f, difference=%0.3f\n",
				status, ototal, oActualTotal, oActualTotal-ototal,
			)
		}

		gtotal += ototal
		gActualTotal += oActualTotal
		fmt.Println()
	}

	var status string
	if gtotal <= gActualTotal {
		status = "PROFIT"
	} else {
		status = "LOSS"
	}

	fmt.Printf("############ All Currencies Total #############\n")
	if gtotal > 0 {
		fmt.Printf(
			"%s -> total_spend=%0.3f, current_worth=%0.3f, difference=%0.3f\n",
			status, gtotal, gActualTotal, gActualTotal-gtotal,
		)
	}

}
