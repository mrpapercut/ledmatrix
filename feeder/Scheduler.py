import schedule
import threading
import time

from .Logger import Logger

from .YoutubeChannels import YoutubeChannels
from .SQLite import SQLite

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
