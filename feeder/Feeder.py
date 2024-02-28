import json
import time

from .Scheduler import Scheduler

class Feeder:
    def __init__(self, config):
        self.config = config

    def start_scheduler(self):
        self.scheduler = Scheduler(self.config)
        self.scheduler.start()

        try:
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            self.scheduler.stop()
            self.scheduler.join()


if __name__ == "__main__":
    config_file = './config.json'
    with open(config_file, 'r') as file:
        json_config = json.load(file)

    feeder = Feeder(json_config)
    feeder.start_scheduler()
