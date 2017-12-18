import os
import json
import requests
import pprint
import base64, hashlib, hmac, time
from collections import OrderedDict


PUBLIC_API_KEY = os.environ["BTCMARKETS_API_KEY"]
PRIVATE_API_KEY = base64.b64decode(os.environ["BTCMARKETS_SECRET_KEY_B64"])
BASE_URL = "https://api.btcmarkets.net"
#BASE_URL = "http://localhost:8001"


def make_request(path, method, data=None):

    now_milliseconds = str(int(time.time() * 1000))
    now_milliseconds = "1513504571000"

    if data is not None:
        data_string = json.dumps(data, separators=(',', ':'))
        string_to_sign = "{path}\n{now}\n{data}".format(
            path=path, now=now_milliseconds, data=data_string
        )
    else:
        string_to_sign = "{path}\n{now}\n".format(
            path=path, now=now_milliseconds
        )


    signature = base64.b64encode(hmac.new(PRIVATE_API_KEY, string_to_sign, digestmod=hashlib.sha512).digest())

    print(string_to_sign)
    print(signature)

    header = {
        'accept': 'application/json',
        'Content-Type': 'application/json',
        'User-Agent': 'python api testing',
        'accept-charset': 'utf-8',
        'apikey': PUBLIC_API_KEY,
        'signature': signature,
        'timestamp': now_milliseconds,
    }

    if method == "GET":
        response = requests.get(BASE_URL+path, headers=header)
    else:
        response = requests.post(BASE_URL+path, json=data, headers=header)

    return response


def main2():
    data = OrderedDict([
        ('currency', "AUD"),
        ('instrument', "XRP"),
        ('limit', 100),
        ('since', 1),
    ])

    #r = make_request("/order/trade/history", "POST", data)
    r = make_request("/order/history", "POST", data)

    #r = make_request("/market/ETH/AUD/tick", "GET")
    #r = make_request("/market/ETH/AUD/orderbook", "GET")
    #r = make_request("/market/ETH/AUD/trades", "GET")

    #r = make_request("/account/balance", "GET")

    print(r.status_code)
    print(r.text)
    pprint.pprint(r.json())


def main():
    string_to_sign = """/order/history
1513501739000
{"currency": "AUD", "instrument": "XRP", "limit": 100, "since": 1}"""
    signature = base64.b64encode(hmac.new(PRIVATE_API_KEY, string_to_sign, digestmod=hashlib.sha512).digest())

    print(signature)


if __name__ == "__main__":
    main2()
