import schedule
import threading
import time

from datetime import datetime

from .YoutubeChannels import YoutubeChannels
from .SQLite import SQLite

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
        print(f"{ts()} Starting scheduler")

        self.init_db()

        self.update_youtube_channels()

        self.is_running = True

        schedule.every(1).minutes.do(self.test_scheduler)
        schedule.every(15).minutes.do(self.update_youtube_channels)

        while self.is_running:
            schedule.run_pending()
            time.sleep(1)

    def stop(self):
        self.is_running = False

    def test_scheduler(self):
        print(f"{ts()} Scheduler is running")

    def update_youtube_channels(self):
        print(f"{ts()} Updating Youtube channels")
        yt = YoutubeChannels(self.config, self.db)
        yt.get_videos()
