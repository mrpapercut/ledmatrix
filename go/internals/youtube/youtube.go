package youtube

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/mrpapercut/ledmatrix/internals/config"
	"github.com/mrpapercut/ledmatrix/internals/sqlite"
	"github.com/mrpapercut/ledmatrix/internals/types"
	"github.com/mrpapercut/ledmatrix/internals/utils"
)

type Youtube struct {
	config                *config.Config
	rssUrlPrefix          string
	videoDetailsUrlPrefix string
}

var youtubeLock = &sync.Mutex{}
var youtubeInstance *Youtube

func GetYoutubeInstance(config *config.Config) *Youtube {
	if youtubeInstance == nil {
		youtubeLock.Lock()
		defer youtubeLock.Unlock()

		if youtubeInstance == nil {
			youtubeInstance = &Youtube{
				config:                config,
				rssUrlPrefix:          "https://www.youtube.com/feeds/videos.xml?channel_id=",
				videoDetailsUrlPrefix: "https://youtube.googleapis.com/youtube/v3/videos?part=snippet,contentDetails,liveStreamingDetails&id=",
			}
		}
	}

	return youtubeInstance
}

func (y *Youtube) getRssFeed(channelId string) (*types.Rss, error) {
	response, err := http.Get(fmt.Sprintf("%s%s", y.rssUrlPrefix, channelId))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rss := types.Rss{}

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
			slog.Warn("Error parsing RSS feed:", err)
			continue
		}

		y.ParseVideos(channelName, feed.Entries)
	}
}

func (y *Youtube) ParseVideos(channel string, videos []types.Entry) {
	datestringLayout := "2006-01-02T15:04:05-07:00"
	sql := sqlite.GetSQLiteInstance()

	for _, video := range videos {
		// Trim off "yt:video:"
		video.VideoId = video.VideoId[9:]

		if y.CheckVideoExistsInDb(video.VideoId) {
			continue
		}

		parsedPublishedTime, err := time.Parse(datestringLayout, video.Published)
		if err != nil {
			slog.Warn("Error parsing published date:", err)
		}

		videoDetails, err := y.getVideoDetails(video.VideoId)
		if err != nil {
			slog.Warn("Error parsing video details:", err)
		}

		if videoDetails.Duration > 0 && videoDetails.Duration <= 60 {
			continue
		}

		youtubeVideo := types.YoutubeVideo{
			VideoId:   video.VideoId,
			Published: parsedPublishedTime.Unix(),
			Channel:   channel,
			Title:     video.Title,
			Details: types.YoutubeVideoDetails{
				Duration:     videoDetails.Duration,
				IsLive:       videoDetails.IsLive,
				IsUpcoming:   videoDetails.IsUpcoming,
				UpcomingDate: videoDetails.UpcomingDate,
			},
		}

		sql.InsertYoutubeVideo(youtubeVideo)

		slog.Info("Inserting video. Title: %v, published: %v, duration: %v seconds\n",
			video.Title, parsedPublishedTime.Format("2006-01-02 15:04:05"), youtubeVideo.Details.Duration)

		var message string
		if youtubeVideo.Details.IsUpcoming {
			scheduledDateString := time.Unix(youtubeVideo.Details.UpcomingDate, 0).Format("Monday 02/01 at 15:04")

			message = fmt.Sprintf("%v scheduled %v to start on %v", youtubeVideo.Channel, youtubeVideo.Title, scheduledDateString)
		} else if youtubeVideo.Details.IsLive {
			message = fmt.Sprintf("%v just went live with %v", youtubeVideo.Channel, youtubeVideo.Title)
		} else {
			message = fmt.Sprintf("New video for %v: %v", youtubeVideo.Channel, youtubeVideo.Title)
		}

		feedMessage := types.FeedMessage{
			Timestamp:    youtubeVideo.Published,
			Type:         "youtubeVideo",
			Message:      message,
			ExtraParam:   youtubeVideo.Channel,
			Priority:     1,
			DisplayTimes: 3,
		}

		sql.InsertFeedMessage(feedMessage)

		// If it's an upcoming video, schedule a message for when the video goes live
		if youtubeVideo.Details.IsUpcoming {
			scheduledMessage := fmt.Sprintf("%v just went live with %v", youtubeVideo.Channel, youtubeVideo.Title)

			scheduledFeedMessage := types.FeedMessage{
				Timestamp:    youtubeVideo.Details.UpcomingDate,
				Type:         "youtubeVideo",
				Message:      scheduledMessage,
				ExtraParam:   youtubeVideo.Channel,
				Priority:     1,
				DisplayTimes: 3,
			}

			sql.InsertFeedMessage(scheduledFeedMessage)
		}
	}
}

func (y *Youtube) CheckVideoExistsInDb(videoId string) bool {
	sqlite := sqlite.GetSQLiteInstance()
	return sqlite.VideoExists(videoId)
}

func (y *Youtube) getVideoDetails(videoId string) (types.YoutubeVideoDetails, error) {
	datestringLayout := "2006-01-02T15:04:05Z"
	url := fmt.Sprintf("%s%s&key=%s", y.videoDetailsUrlPrefix, videoId, y.config.Youtube.ApiKey)

	videoDetails := types.YoutubeVideoDetails{}

	response, err := http.Get(url)
	if err != nil {
		return videoDetails, err
	}
	defer response.Body.Close()

	var jsonVideoDetails types.VideoDetails

	err = json.NewDecoder(response.Body).Decode(&jsonVideoDetails)
	if err != nil {
		return videoDetails, err
	}

	if len(jsonVideoDetails.Items) < 1 {
		slog.Warn("No videodetails found for %v: %v\n", videoId, jsonVideoDetails)
		return videoDetails, nil
	}

	videoDetails.IsLive = jsonVideoDetails.Items[0].Snippet.Live == "live"
	videoDetails.IsUpcoming = jsonVideoDetails.Items[0].Snippet.Live == "upcoming"

	if len(jsonVideoDetails.Items[0].LiveStreamingDetails.ScheduledStartTime) > 0 {
		parsedScheduledTime, err := time.Parse(datestringLayout, jsonVideoDetails.Items[0].LiveStreamingDetails.ScheduledStartTime)
		if err != nil {
			slog.Error("Error parsing scheduled date:", err)
		}
		videoDetails.UpcomingDate = parsedScheduledTime.Unix()
	}

	duration, err2 := utils.ParseDurationString(jsonVideoDetails.Items[0].Details.Duration)
	if err2 != nil {
		return videoDetails, err2
	}
	videoDetails.Duration = duration

	return videoDetails, nil
}
