import logging

class Logger:
    def __init__(self, config) -> None:
        self.log_format = '%(asctime)s [%(name)s::%(levelname)s] %(message)s'
        self.date_format = '%Y-%m-%d %H:%M:%S'

        self.log_to_file = config.get('log').get('log_to_file')
        self.log_file = config.get('log').get('filename')

        self.log_level = {
            'debug': logging.DEBUG,
            'info': logging.INFO,
            'warning': logging.WARNING,
            'error': logging.ERROR,
            'critical': logging.CRITICAL
        }[config.get('log').get('level')]

        self.setup_logger()

    def setup_logger(self):
        logger = logging.getLogger()

        logger.setLevel(self.log_level)

        formatter = logging.Formatter(self.log_format, datefmt=self.date_format)

        if self.log_to_file:
            file_handler = logging.FileHandler(self.log_file, encoding='utf-8')
            file_handler.setFormatter(formatter)
            logger.addHandler(file_handler)
        else:
            console_handler = logging.StreamHandler()
            console_handler.setFormatter(formatter)
            logger.addHandler(console_handler)

    def get_logger(self, name):
        return logging.getLogger(name)
