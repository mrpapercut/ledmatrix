import sqlite3

class SQLite:
    def __init__(self, config) -> None:
        self.dbfilename = config.get('db').get('filename')

    def connect(self):
        self.db = sqlite3.connect(self.dbfilename)

    def close(self):
        self.db.close()

    def init_db(self):
        cursor = self.db.cursor()

        # cursor.execute("CREATE TABLE IF NOT EXISTS feed(id, timestamp, )")
