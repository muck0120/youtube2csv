package fixture

import (
	"time"

	"github.com/muck0120/youtube2csv/internal/domain/youtube"
)

func Chnanel(setter func(channel *youtube.Channel)) *youtube.Channel {
	channel := &youtube.Channel{
		ID: "ABCDEFG123456",
		Videos: []*youtube.Video{
			Video(func(video *youtube.Video) {
				video.ID = "AAAAA"
				video.Title = "video1"
			}),
			Video(func(video *youtube.Video) {
				video.ID = "BBBBB"
				video.Title = "video2"
			}),
			Video(func(video *youtube.Video) {
				video.ID = "CCCCC"
				video.Title = "video3"
			}),
		},
	}

	if setter != nil {
		setter(channel)
	}

	return channel
}

func Video(setter func(video *youtube.Video)) *youtube.Video {
	video := &youtube.Video{
		ID:       "AAAAA",
		Title:    "video1",
		Duration: 10 * time.Minute,
		Status:   youtube.VideoStatusPublic,
	}

	if setter != nil {
		setter(video)
	}

	return video
}
