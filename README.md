# BTC Markets Golang API
This is a basic implementation of the BTC Markets api. It uses your api keys to access your accounts current cyrptocurrency balance. It analyses your transactions and determines whether you are operating at a net profit, and which trades have made you a loss or profit.

## Usage
```
NAME:
   BTC Markets - Account and market prices for BTC Markets cryptocurrencies

USAGE:
   BTC Markets [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
    account	Account balance and analysis

GLOBAL OPTIONS:
   --public_key 	Public api key [$BTC_MARKETS_PUBLIC_API_KEY]
   --private_key 	Base64 encoded private key [$BTC_MARKETS_PRIVATE_API_KEY]
   --help, -h		show help
   --version		print the version
```

### Account command
```
$ go run main.go --public_key=<public_key> --private_key=<private_key> account
```

#### Example output
```
############ AUD #############
Balance 105.449AUD

############ BTC #############
Balance 0.000BTC
Current market price: last_price=26180.000, best_bid=26150.000, best_ask=26180.000

############ LTC #############
Balance 5.000LTC
Current market price: last_price=439.600, best_bid=439.600, best_ask=439.800
Transactions:
>> 437.970AUD * 0.010 = 4.380AUD 	| current value: 4.396AUD 	| PROFIT 0.016
>> 440.000AUD * 0.500 = 220.000AUD 	| current value: 219.800AUD 	| LOSS -0.200
PROFIT -> total_spend=1023.350, current_worth=1103.396, difference=80.046

############ All Currencies Total #############
PROFIT -> total_spend=1234.54, current_worth=1234.54, difference=111
```

## TODO
```
- Make the account analyzer calls aysnc
- Add command to watch all or specific currencies live
```
