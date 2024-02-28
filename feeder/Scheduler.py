import schedule
import threading
import time
import logging

from datetime import datetime

from .YoutubeChannels import YoutubeChannels
from .SQLite import SQLite

logging.basicConfig(level=logging.INFO)

def ts():
    return f"[{datetime.now().strftime('%H:%M:%S')}]"

class Scheduler(threading.Thread):
    def __init__(self, config):
        super().__init__()

        self.config = config

        self.db = SQLite(self.config)

        self.is_running = False

    def init_db(self):
        if not self.db.connected():
            self.db.connect()

        self.db.init_db()

    def run(self):
        logging.info(f"{ts()} Starting scheduler")

        self.init_db()

        self.is_running = True

        schedule.every(1).minutes.do(self.test_scheduler)
        schedule.every(15).minutes.do(self.update_youtube_channels)

        while self.is_running:
            schedule.run_pending()
            time.sleep(1)

        logging.info(f"{ts()} Scheduler shutting down")
        self.db.close()

    def stop(self):
        self.is_running = False

    def test_scheduler(self):
        logging.info(f"{ts()} Scheduler is running")

    def update_youtube_channels(self):
        logging.info(f"{ts()} Updating Youtube channels")

        yt = YoutubeChannels(self.config, self.db)
        yt.get_videos()

        logging.info(f"{ts()} Finished updating Youtube channels")
