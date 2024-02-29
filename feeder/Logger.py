import logging

class Logger:
    def __init__(self, log_file='feeder.log', log_level=logging.DEBUG) -> None:
        self.log_format = '%(asctime)s [%(name)s::%(levelname)s] %(message)s'
        self.date_format = '%Y-%m-%d %H:%M:%S'

        self.log_file = log_file

        self.log_level = log_level

        self.setup_logger()

    def setup_logger(self):
        logging.basicConfig(level=self.log_level, format=self.log_format, datefmt=self.date_format)

    def get_logger(self, name):
        return logging.getLogger(name)
