package youtube

import (
	"context"

	"muck0120/youtube2csv/internal/domain/youtube"
	"muck0120/youtube2csv/pkg/errors"
	"muck0120/youtube2csv/pkg/time"
)

type Repository struct {
	Service IService
}

func NewRepository(service IService) *Repository {
	return &Repository{Service: service}
}

func (rp *Repository) FindByID(_ context.Context, channelID string) (*youtube.Channel, error) {
	youtubeVideos, err := rp.Service.FindVideosByChannelID(channelID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	videos := make([]*youtube.Video, 0, len(youtubeVideos))

	for _, youtubeVideo := range youtubeVideos {
		duration, err := time.ParseISO8601Duration(youtubeVideo.ContentDetails.Duration)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		videos = append(videos, &youtube.Video{
			ID:       youtubeVideo.Id,
			Title:    youtubeVideo.Snippet.Title,
			Duration: duration,
			Status:   youtube.VideoStatus(youtubeVideo.Status.PrivacyStatus),
		})
	}

	return &youtube.Channel{
		ID:     channelID,
		Videos: videos,
	}, nil
}
