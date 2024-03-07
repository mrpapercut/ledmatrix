import schedule
import threading
import time

from .Logger import Logger

from .SQLite import SQLite
from .Weather import Weather
from .YoutubeChannels import YoutubeChannels

class Scheduler(threading.Thread):
    def __init__(self, config):
        super().__init__()

        self.config = config

        self.logger = Logger(config).get_logger(__name__)

        self.db = SQLite(self.config)

        self.is_running = False

    def init_db(self):
        if not self.db.connected():
            self.db.connect()

        self.db.init_db()

    def run(self):
        self.logger.info("Starting scheduler")

        self.init_db()

        self.is_running = True

        schedule.every(1).minutes.do(self.test_scheduler)
        schedule.every(15).minutes.do(self.update_youtube_channels)
        schedule.every(1).hours.do(self.update_current_weather)

        while self.is_running:
            schedule.run_pending()
            time.sleep(1)

        self.logger.info("Scheduler shutting down")
        self.db.close()

    def stop(self):
        self.is_running = False

    def test_scheduler(self):
        self.logger.info("Scheduler is running")

    def update_youtube_channels(self):
        self.logger.info("Updating Youtube channels")

        yt = YoutubeChannels(self.config, self.db)
        yt.get_videos()

        self.logger.info("Finished updating Youtube channels")

    def update_current_weather(self):
        self.logger.info("Updating current weather")

        weather = Weather(self.config, self.db)
        weather.get_current_conditions()

        self.logger.info("Finished updating current weather")
