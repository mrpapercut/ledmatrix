import feedparser
import requests
import time

from isodate import parse_duration

class YoutubeChannels:
    def __init__(self, config, db) -> None:
        self.db = db

        self.youtube_api_key = config.get("youtube").get("apikey")
        self.channels = config.get("youtube").get("channels")

        self.youtube_rss_url_prefix = "https://www.youtube.com/feeds/videos.xml?channel_id="
        self.youtube_video_details_url_prefix = "https://youtube.googleapis.com/youtube/v3/videos?part=contentDetails&id="

    def get_videos(self) -> None:
        videos = {}

        for channel, channel_id in self.channels.items():
            feed = feedparser.parse(f"{self.youtube_rss_url_prefix}{channel_id}")

            videos[channel] = self.parse_rss_entries(feed.entries)
            self.add_durations_to_videos(videos[channel])

        self.db.insert_youtube_videos(videos)

    def parse_rss_entries(self, entries) -> list:
        videos = []

        for entry in entries:
            video_id = entry.yt_videoid
            title = entry.title
            published = int(time.mktime(entry.published_parsed))
            updated = int(time.mktime(entry.updated_parsed))

            videos.append({
                "id": video_id,
                "title": title,
                "published": published,
                "updated": updated,
            })

        return videos

    def add_durations_to_videos(self, videos):
        video_ids = [video['id'] for video in videos]

        r = requests.get(f"{self.youtube_video_details_url_prefix}{','.join(video_ids)}&key={self.youtube_api_key}")
        json_response = r.json()

        durations = {response['id']: parse_duration(response['contentDetails']['duration']).total_seconds() for response in json_response['items']}

        for video in videos:
            video_id = video['id']
            if video_id in durations:
                video['duration'] = durations[video_id]
