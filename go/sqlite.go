package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type FeedMessage struct {
	Timestamp    int64
	Type         string
	Message      string
	ExtraParam   string
	Priority     int
	DisplayTimes int
}

type SQLite struct {
	config *Config
}

var sqliteLock = &sync.Mutex{}
var sqliteInstance *SQLite

func getSQLiteInstance() *SQLite {
	if sqliteInstance == nil {
		sqliteLock.Lock()
		defer sqliteLock.Unlock()

		if sqliteInstance == nil {
			sqliteInstance = &SQLite{}
			sqliteInstance.init()
		}
	}

	return sqliteInstance
}

func (s *SQLite) init() {
	s.config = getConfig()

	if _, err := os.Stat(s.config.Database.Filename); os.IsNotExist(err) {
		// If it doesn't exist, create it
		_, err := os.Create(s.config.Database.Filename)
		if err != nil {
			fmt.Println("Error creating database file:", err)
			return
		}
	}

	db, err := sql.Open("sqlite3", s.config.Database.Filename)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	createTablesQueries := []string{
		`CREATE TABLE IF NOT EXISTS feed (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp INTEGER NOT NULL,
			type TEXT NOT NULL,
			message TEXT NOT NULL,
			param TEXT NOT NULL,
			priority INTEGER,
			display_times INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS yt_videos (
			video_id TEXT NOT NULL,
			published INTEGER NOT NULL,
			channel TEXT NOT NULL,
			title TEXT NOT NULL,
			duration INTEGER,
			UNIQUE (video_id)
		);`,
	}

	for _, query := range createTablesQueries {
		_, err = db.Exec(query)
		if err != nil {
			fmt.Println("Error creating tables:", err)
		}
	}
}

func (s *SQLite) InsertFeedMessage(feedMessage FeedMessage) {
	db, err := sql.Open("sqlite3", s.config.Database.Filename)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO feed(timestamp, type, message, param, priority, display_times) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error creating statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		feedMessage.Timestamp,
		feedMessage.Type,
		feedMessage.Message,
		feedMessage.ExtraParam,
		feedMessage.Priority,
		feedMessage.DisplayTimes,
	)
	if err != nil {
		fmt.Println("Error inserting feed message:", err)
	}
}

func (s *SQLite) GetLatestFeedMessage(priority int) (FeedMessage, error) {
	var result FeedMessage

	db, err := sql.Open("sqlite3", s.config.Database.Filename)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return result, err
	}
	defer db.Close()

	getFeedMessageQuery := `SELECT timestamp, type, message, param, priority, display_times
	FROM feed
	WHERE priority = ?
	AND display_times > 0
	AND timestamp > ?
	ORDER BY display_times DESC, timestamp DESC LIMIT 1`
	yesterday := time.Now().Add(-24 * time.Hour).Unix()
	err = db.QueryRow(getFeedMessageQuery, priority, yesterday).Scan(&result.Timestamp, &result.Type, &result.Message, &result.ExtraParam, &result.Priority, &result.DisplayTimes)

	if err == sql.ErrNoRows {
		fmt.Println("No rows found")
		return result, nil
	} else if err != nil {
		fmt.Println("Error retrieving data from db:", err)
		return result, err
	}

	return result, nil
}

func (s *SQLite) LowerPriorityAndDisplayTimes(message FeedMessage) {
	db, err := sql.Open("sqlite3", s.config.Database.Filename)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE feed SET display_times = display_times - 1, priority = 2 WHERE message = ?")
	if err != nil {
		fmt.Println("Error creating statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.Message)
	if err != nil {
		fmt.Printf("Error lowering display_times for %v: %v\n", message.Message, err)
	}
}

func (s *SQLite) InsertYoutubeVideo(video YoutubeVideo) {
	db, err := sql.Open("sqlite3", s.config.Database.Filename)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO yt_videos(video_id, published, channel, title, duration) values(?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error creating statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(video.VideoId, video.Published, video.Channel, video.Title, video.Duration)
	if err != nil {
		fmt.Println("Error inserting youtube video:", err)
	}
}

func (s *SQLite) VideoExists(videoId string) bool {
	db, err := sql.Open("sqlite3", s.config.Database.Filename)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT video_id FROM yt_videos WHERE video_id = ?")
	if err != nil {
		fmt.Println("Error creating statement:", err)
		return false
	}
	defer stmt.Close()

	var dbVideoId string

	err = stmt.QueryRow(videoId).Scan(&dbVideoId)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		fmt.Println("Error retrieving data from db:", err)
	}

	return true
}
