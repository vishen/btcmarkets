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


## TODO
```
- Make the account analyzer calls aysnc
- Add command to watch all or specific currencies live
```
