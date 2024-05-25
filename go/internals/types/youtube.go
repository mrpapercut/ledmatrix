package types

type YoutubeVideo struct {
	VideoId   string
	Published int64
	Channel   string
	Title     string
	Details   YoutubeVideoDetails
}

type YoutubeVideoDetails struct {
	Duration     int64
	IsLive       bool
	IsUpcoming   bool
	UpcomingDate int64
}

type Entry struct {
	VideoId   string `xml:"id"`
	Published string `xml:"published"`
	Title     string `xml:"title"`
}

type Rss struct {
	Entries []Entry `xml:"entry"`
}

type VideoDetailsItem struct {
	Id      string `json:"id"`
	Details struct {
		Duration string `json:"duration"`
	} `json:"contentDetails"`
	Snippet struct {
		Live string `json:"liveBroadcastContent"`
	} `json:"snippet"`
	LiveStreamingDetails struct {
		ScheduledStartTime string `json:"scheduledStartTime"`
	} `json:"liveStreamingDetails"`
}

type VideoDetails struct {
	Items []VideoDetailsItem `json:"items"`
}
