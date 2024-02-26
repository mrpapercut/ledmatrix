import json
import schedule

from .SQLite import SQLite
from .YoutubeChannels import YoutubeChannels

class Feeder:
    def __init__(self, config):
        self.config = config

    def start_schedule(self):
        self.schedule = SQLite(self.config)

    def update_youtube_channels(self):
        yt = YoutubeChannels(self.config)
        yt.get_videos()


if __name__ == "__main__":
    config_file = './config.json'
    with open(config_file, 'r') as file:
        json_config = json.load(file)

    feeder = Feeder(json_config)
    feeder.start_schedule()
    feeder.update_youtube_channels()
