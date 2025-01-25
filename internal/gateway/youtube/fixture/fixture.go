package fixture

import (
	"google.golang.org/api/youtube/v3"
)

func Videos(setter func(videos []*youtube.Video)) []*youtube.Video {
	videos := []*youtube.Video{
		Video(func(video *youtube.Video) {
			video.Id = "AAAAA"
			video.ContentDetails.Duration = "PT1H1M1S"
			video.Snippet.Title = "video1"
		}),
		Video(func(video *youtube.Video) {
			video.Id = "BBBBB"
			video.ContentDetails.Duration = "PT1H2M2S"
			video.Snippet.Title = "video2"
		}),
		Video(func(video *youtube.Video) {
			video.Id = "CCCCC"
			video.ContentDetails.Duration = "PT1H3M3S"
			video.Snippet.Title = "video3"
		}),
	}

	if setter != nil {
		setter(videos)
	}

	return videos
}

func Video(setter func(video *youtube.Video)) *youtube.Video {
	video := &youtube.Video{
		Id:             "ABCDEFG123456",
		ContentDetails: &youtube.VideoContentDetails{Duration: "PT1H40M"},
		Snippet:        &youtube.VideoSnippet{Title: "ABCDE"},
		Status:         &youtube.VideoStatus{PrivacyStatus: "public"},
	}

	if setter != nil {
		setter(video)
	}

	return video
}
