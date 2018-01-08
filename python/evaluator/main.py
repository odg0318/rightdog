import sys
import os
import yaml
import time

from influxdb import InfluxDBClient
from slacker import Slacker


config = {}
slack = None


def main():
    global config
    global slack

    config_path = os.getenv('EVALUATOR_CONFIG', 'evaluator.yml')
    config = yaml.load(file(config_path, 'r'))

    slack = Slacker(config['slack']['token'])
    exchanges = config['exchange']['ko']

    while True:
        process_korean_arbitrage(exchanges)
        time.sleep(config['trade']['interval'])


def process_korean_arbitrage(exchanges):
    for exchange in exchanges:
        data = make_data_from_db(exchange)

        for currency, prices in data.items():
            if len(prices) == 1:
                continue

            base_price = prices[exchange]

            for e in prices.keys():
                if e == exchange:
                    continue

                gap_rate = caculate_gap_rate(base_price, prices[e])
                if gap_rate > 3.:
                    message = 'base / %s / %d\ntarget / %s / %d\ngap / %f' % (exchange, base_price, e, prices[e], gap_rate)
                    slack.chat.post_message(config['slack']['channel'], message)


def make_data_from_db(exchange):
    client = InfluxDBClient(host='127.0.0.1', port=8086, database='rightdog')

    query = "SHOW TAG VALUES FROM ticker WITH KEY = fromcurrency WHERE exchange = '%s' AND time > now() - 10m" % exchange
    result = client.query(query)

    data = {}
    for row in result.get_points():
        currency = row['value']

        query = "SELECT LAST(price) FROM ticker WHERE fromcurrency = '%s' AND exchange =~ /%s/ GROUP BY exchange" % (currency, '|'.join(exchanges))
        result = client.query(query)

        data[currency] = {}
        for key in result.keys():
            data[currency][key[1]['exchange']] = next(result[key])['last']

    # {u'bch': {u'coinone': 4065500, u'upbit': 4085000, u'korbit': 25950000}, u'etc': {u'coinone': 53960, u'upbit': 54810, u'korbit': 25954000}, u'btg': {u'coinone': 589100, u'upbit': 399050, u'korbit': 25954000}, u'btc': {u'coinone': 25350000, u'upbit': 26149000, u'korbit': 25950000}, u'eth': {u'coinone': 1673100, u'upbit': 1692500, u'korbit': 25954000}, u'xrp': {u'coinone': 4400, u'upbit': 4330, u'korbit': 25950000}}
    return data


def caculate_gap_rate(base, target):
    return (float(target)/base - 1.0) * 100


if __name__ == '__main__':
    main()
