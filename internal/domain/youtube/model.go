package youtube

import "time"

type Channel struct {
	ID     string
	Videos []*Video
}

type Video struct {
	ID       string
	Title    string
	Duration time.Duration
	Status   VideoStatus
}

type VideoStatus string

const (
	VideoStatusPrivate  VideoStatus = "private"
	VideoStatusPublic   VideoStatus = "public"
	VideoStatusUnlisted VideoStatus = "unlisted"
)
