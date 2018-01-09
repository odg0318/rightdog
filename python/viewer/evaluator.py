from tornado.options import options, define
from influxdb import InfluxDBClient

def connect_influxdb():
    return InfluxDBClient(host=options.config['influxdb']['host'], port=options.config['influxdb']['port'], database=options.config['influxdb']['database'])


def make_data_from_db(exchanges, currencies):
    client = connect_influxdb()

    data = {}
    for currency in currencies:
        query = "SELECT LAST(price) FROM ticker WHERE fromcurrency = '%s' AND exchange =~ /%s/  GROUP BY exchange" % (currency, '|'.join(exchanges))
        result = client.query(query)

        data[currency] = {}
        for key in result.keys():
            exchange = key[1]['exchange'].encode('utf-8')
            price = next(result[key])['last']

            data[currency][exchange] = {
                'price': price
            }

        min_price = min([y['price'] for x, y in data[currency].iteritems()])
        max_price = max([y['price'] for x, y in data[currency].iteritems()])

        for exchange in exchanges:
            data[currency][exchange]['gap_rate'] = calculate_gap_rate(min_price, data[currency][exchange]['price'])
            data[currency][exchange]['is_min'] = min_price == data[currency][exchange]['price']
            data[currency][exchange]['is_max'] = max_price == data[currency][exchange]['price']

    return data


def calculate_gap_rate(base, target):
    return (float(target)/base - 1.0) * 100


if __name__ == '__main__':
    define('config', default={})

    options.config = {
        'influxdb': {
            'host': '127.0.0.1',
            'port': 8086,
            'database': 'rightdog'
        }
    }
    data = make_data_from_db(['korbit', 'coinone', 'upbit'], ['btc', 'eth', 'etc'])

    print data
