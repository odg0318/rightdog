import yaml

import tornado.ioloop
import tornado.web
from tornado.options import define, options, parse_command_line

import evaluator


class MainHandler(tornado.web.RequestHandler):
    def get(self):
        data = evaluator.make_data_from_db(options.config['exchange']['ko'], options.config['currency']['ko'])
        self.render('views/main.html', data=data, exchanges=options.config['exchange']['ko'], currencies=options.config['currency']['ko'])


def make_app():
    define('config_path', default='./viewer.yml', help='Configuration path')
    define('config', default={})

    parse_command_line()

    options.config = yaml.load(file(options.config_path, 'r'))

    return tornado.web.Application([
        (r"/", MainHandler),
    ])


if __name__ == '__main__':
    app = make_app()
    app.listen(20000)
    tornado.ioloop.IOLoop.current().start()
