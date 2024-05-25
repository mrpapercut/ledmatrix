package types

type FeedMessage struct {
	Timestamp    int64
	Type         string
	Message      string
	ExtraParam   string
	Priority     int
	DisplayTimes int
}
