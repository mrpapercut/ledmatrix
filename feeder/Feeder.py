import json
import time

from .Scheduler import Scheduler
from .SQLite import SQLite

class Feeder:
    def __init__(self, config):
        self.config = config
        self.db = SQLite(self.config)
        self.init_db()

    def init_db(self):
        if not self.db.connected():
            self.db.connect()

        self.db.init_db()

    def start_scheduler(self):
        self.scheduler = Scheduler(self.config, self.db)
        self.scheduler.start()

        try:
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            self.scheduler.join()



if __name__ == "__main__":
    config_file = './config.json'
    with open(config_file, 'r') as file:
        json_config = json.load(file)

    feeder = Feeder(json_config)
    feeder.start_scheduler()
