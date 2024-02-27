import sqlite3

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
                message TEXT NOT NULL
            )
        """)
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS yt_videos(
                video_id TEXT,
                UNIQUE (video_id)
            )
        """)

        self.db.commit()

        cursor.close()

    def insert_feed_message(self, timestamp, type, message):
        cursor = self.db.cursor()

        try:
            cursor.execute("INSERT INTO feed (timestamp, type, message) VALUES (?, ?, ?)", (timestamp, type, message))

            self.db.commit()
        except Exception as e:
            print(f"Error inserting feed message: {e}")

        cursor.close()

    def insert_youtube_videos(self, videos):
        cursor = self.db.cursor()

        for channel, videos in videos.items():
            for video in videos:
                try:
                    video_id = video.get('id')
                    cursor.execute("INSERT INTO yt_videos (video_id) VALUES (?)", (video_id,))

                    self.db.commit()

                    self.insert_feed_message(video.get('published'), ':yt:newvideo:', f"New video for {channel}: {video.get('title')}")
                except sqlite3.IntegrityError:
                    # video already exists
                    pass

        cursor.close()
