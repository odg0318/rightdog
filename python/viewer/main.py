import os
import yaml

import tornado.ioloop
import tornado.web
from tornado.options import define, options, parse_command_line

from api import ApiGetPriceHandler

import evaluator

class MainHandler(tornado.web.RequestHandler):
    def get(self):
        self.render('views/main.html')


def make_app():
    define('config_path', default='./viewer.yml', help='Configuration path')
    define('config', default={})

    parse_command_line()

    options.config = yaml.load(file(options.config_path, 'r'))

    return tornado.web.Application([
        (r"/", MainHandler),
        (r'/static/(.*)', tornado.web.StaticFileHandler, {'path': os.path.join(os.path.dirname(__file__), 'static')}),
        (r"/api/price", ApiGetPriceHandler),
    ], autoreload=True)


if __name__ == '__main__':
    app = make_app()
    app.listen(20000)
    tornado.ioloop.IOLoop.current().start()
