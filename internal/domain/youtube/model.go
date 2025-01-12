package youtube

import (
	"regexp"
	"strconv"
	"time"
)

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

func (v *Video) GetNumber() int {
	matches := regexp.MustCompile(`#(\d+)`).FindStringSubmatch(v.Title)
	if len(matches) < (1 + 1) { // マッチした文字列が1つ以上あるか
		return 0
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}

	return num
}

func (v *Video) GetURL() string {
	return "https://youtu.be/" + v.ID
}

type VideoStatus string

const (
	VideoStatusPrivate  VideoStatus = "private"
	VideoStatusPublic   VideoStatus = "public"
	VideoStatusUnlisted VideoStatus = "unlisted"
)
