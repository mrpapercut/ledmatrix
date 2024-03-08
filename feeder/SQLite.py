import sqlite3
from datetime import datetime

class SQLite:
    def __init__(self, config) -> None:
        self.dbfilename = config.get('db').get('filename')
        self.db = None

    def connect(self) -> None:
        self.db = sqlite3.connect(self.dbfilename)

    def connected(self) -> bool:
        return self.db is not None

    def close(self) -> None:
        self.db.close()

    def init_db(self) -> None:
        cursor = self.db.cursor()

        cursor.execute("""
            CREATE TABLE IF NOT EXISTS feed(
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                timestamp TEXT NOT NULL,
                type TEXT NOT NULL,
                message TEXT NOT NULL,
                param TEXT NOT NULL,
                priority INTEGER,
                display_times INTEGER,
            )
        """)
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS yt_videos(
                video_id TEXT,
                duration INTEGER,
                UNIQUE (video_id)
            )
        """)

        self.db.commit()

        cursor.close()

    def insert_feed_message(self, timestamp, type, message, extraParam, priority = 10, displayTimes = 1):
        cursor = self.db.cursor()

        insert_query = "INSERT INTO feed (timestamp, type, message, param, priority, display_times) VALUES (?, ?, ?, ?, ?, ?)"

        try:
            cursor.execute(insert_query, (timestamp, type, message, extraParam, priority, displayTimes))

            self.db.commit()

            print(f"Inserted new message: {message}")
        except Exception as e:
            print(f"Error inserting feed message: {e}")

        cursor.close()

    def insert_youtube_videos(self, videos):
        cursor = self.db.cursor()

        for channel, videos in videos.items():
            for video in videos:
                try:
                    video_id = video.get('id')
                    duration = video.get('duration')

                    # Skip YTShorts
                    if duration < 61:
                        continue

                    cursor.execute("INSERT INTO yt_videos (video_id, duration) VALUES (?, ?)", (video_id, duration))

                    self.db.commit()

                    self.insert_feed_message(video.get('published'), ':yt:newvideo:', f"New video for {channel}: {video.get('title')} ({duration})", channel)
                except sqlite3.IntegrityError:
                    # video already exists in db
                    pass

        cursor.close()

    def insert_current_weather(self, weather):
        cursor = self.db.cursor()

        message = f"Temperature: {weather.get('temperature')} - real feel: {weather.get('real_feel_temperature')} - precipitation: {weather.get('precipitation')}"

        self.insert_feed_message(weather.get('timestamp'), ':weather:current_conditions:', message, "")

        cursor.close()
