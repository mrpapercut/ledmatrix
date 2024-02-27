import schedule
import threading
import time

from datetime import datetime

from .YoutubeChannels import YoutubeChannels

def ts():
    return f"[{datetime.now().strftime('%H:%M:%S')}]"

class Scheduler(threading.Thread):
    def __init__(self, config, db):
        super().__init__()

        self.config = config
        self.db = db

    def run(self):
        print(f"{ts()} Starting scheduler")
        schedule.every(1).minutes.do(self.test_scheduler)
        schedule.every(15).minutes.do(self.update_youtube_channels)

        while True:
            schedule.run_pending()
            time.sleep(1)

    def test_scheduler(self):
        print(f"{ts()} Scheduler is running")

    def update_youtube_channels(self):
        print(f"{ts()} Updating Youtube channels")
        yt = YoutubeChannels(self.config, self.db)
        yt.get_videos()
