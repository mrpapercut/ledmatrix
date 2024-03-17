package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type YoutubeVideo struct {
	VideoId string
	Published int64
	Channel string
	Title string
	Duration int64
}

type Entry struct {
	VideoId string `xml:"id"`
	Published string `xml:"published"`
	Title string `xml:"title"`
}

type Rss struct {
	Entries []Entry `xml:"entry"`
}

type VideoDetailsItem struct {
	Id string `json:"id"`
	Details struct {
		Duration string `json:"duration"`
	} `json:"contentDetails"`
}

type VideoDetails struct {
	Items []VideoDetailsItem `json:"items"`
}

type Youtube struct {
	config *Config
	rssUrlPrefix string
	videoDetailsUrlPrefix string
}

var youtubeLock = &sync.Mutex{}
var youtubeInstance *Youtube

func getYoutubeInstance(config *Config) *Youtube {
	if youtubeInstance == nil {
		youtubeLock.Lock()
		defer youtubeLock.Unlock()

		if youtubeInstance == nil {
			youtubeInstance = &Youtube{
				config: config,
				rssUrlPrefix: "https://www.youtube.com/feeds/videos.xml?channel_id=",
				videoDetailsUrlPrefix: "https://youtube.googleapis.com/youtube/v3/videos?part=contentDetails&id=",
			}
		}
	}

	return youtubeInstance
}

func (y *Youtube) getRssFeed(channelId string) (*Rss, error) {
	response, err := http.Get(fmt.Sprintf("%s%s", y.rssUrlPrefix, channelId))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rss := Rss{}

	decoder := xml.NewDecoder(response.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		return nil, err
	}

	return &rss, nil
}

func (y *Youtube) GetVideos() {
	for channelName, channelId := range y.config.Youtube.Channels {
		feed, err := y.getRssFeed(channelId)

		if err != nil {
			fmt.Println("Error parsing RSS feed:", err)
			continue
		}

		y.ParseVideos(channelName, feed.Entries)
	}
}

func (y *Youtube) ParseVideos(channel string, videos []Entry) {
	datestringLayout := "2006-01-02T15:04:05-07:00"
	sqlite := getSQLiteInstance()

	for _, video := range videos {
		// Trim off "yt:video:"
		video.VideoId = video.VideoId[9:]

		if y.CheckVideoExistsInDb(video.VideoId) {
			continue
		}

		parsedPublishedTime, err := time.Parse(datestringLayout, video.Published)
		if err != nil {
			fmt.Println("Error parsing published date:", err)
		}

		duration, err := y.getVideoDuration(video.VideoId)
		if err != nil {
			fmt.Println("Error parsing json response:", err)
		}

		if duration > 0 && duration <= 60 {
			continue
		}

		youtubeVideo := YoutubeVideo{
			VideoId: video.VideoId,
			Published: parsedPublishedTime.Unix(),
			Channel: channel,
			Title: video.Title,
			Duration: duration,
		}

		sqlite.InsertYoutubeVideo(youtubeVideo)

		fmt.Printf("Inserting video. Title: %v, published: %v, duration: %v seconds\n",
		video.Title, parsedPublishedTime.Format("2006-02-01 15:04:05"), duration)

		// messageType string, message string, param string, priority int, displayTimes int
		message := fmt.Sprintf("New video for %v: %v", youtubeVideo.Channel, youtubeVideo.Title)
		feedMessage := FeedMessage{
			Timestamp: youtubeVideo.Published,
			Type: "youtubeVideo",
			Message: message,
			ExtraParam: youtubeVideo.Channel,
			Priority: 1,
			DisplayTimes: 3,
		}
		sqlite.InsertFeedMessage(feedMessage)
	}
}

func (y *Youtube) CheckVideoExistsInDb(videoId string) bool {
	sqlite := getSQLiteInstance()
	return sqlite.VideoExists(videoId)
}

func (y *Youtube) getVideoDuration(videoId string) (int64, error) {
	url := fmt.Sprintf("%s%s&key=%s", y.videoDetailsUrlPrefix, videoId, y.config.Youtube.ApiKey)

	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	var videoDetails VideoDetails

	err = json.NewDecoder(response.Body).Decode(&videoDetails)
	if err != nil {
		return 0, err
	}

	if len(videoDetails.Items) < 1 {
		fmt.Printf("No videodetails found for %v: %v\n", videoId, videoDetails)
		return 0, nil
	}

	duration, err2 := ParseDurationString(videoDetails.Items[0].Details.Duration)
	if err2 != nil {
		return 0, err2
	}

	return duration, nil
}
