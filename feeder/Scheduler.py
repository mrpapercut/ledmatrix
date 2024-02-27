import schedule
import threading
import time

from .YoutubeChannels import YoutubeChannels

class Scheduler(threading.Thread):
    def __init__(self, config, db):
        super().__init__()

        self.config = config
        self.db = db

        self.start_schedule()

    def start_schedule(self):
        print("Starting scheduler")
        schedule.every(10).minutes.do(self.update_youtube_channels)
        schedule.every(1).minutes.do(self.test_scheduler)

        while True:
            schedule.run_pending()
            time.sleep(1)

    def test_scheduler(self):
        print("Scheduler is running")

    def update_youtube_channels(self):
        print("Updating Youtube channels")
        yt = YoutubeChannels(self.config, self.db)
        yt.get_videos()
