import json

import tornado.web

import evaluator

class BaseApiHandler(tornado.web.RequestHandler):

    def set_default_headers(self):
        self.set_header('Content-Type', 'application/json')


class ApiGetPriceHandler(BaseApiHandler):

    def get(self):
        exchanges = self.get_argument('exchanges').split(',')
        currencies = self.get_argument('currencies').split(',')
        data = {
            'exchanges': exchanges,
            'currencies': currencies,
            'data': evaluator.make_data_from_db(exchanges, currencies)
        }

        self.write(data)
