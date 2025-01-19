//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=$GOPACKAGE -destination=./get_info_mock.go

package youtube

import (
	"context"
	"math"
	"sort"

	"github.com/muck0120/youtube2csv/internal/domain/youtube"
	"github.com/muck0120/youtube2csv/internal/pkg/errors"
)

type IGetInfoUseCase interface {
	Execute(ctx context.Context, in *GetInfoUseCaseInput) (*GetInfoUseCaseOutput, error)
}

type GetInfoUseCase struct {
	youtubeRepository youtube.IRepository
}

type GetInfoUseCaseInput struct {
	ChannelID string
}

type GetInfoUseCaseOutput struct {
	Videos []*Video
}

type Video struct {
	Number         int
	Title          string
	MinuteDuration int
	URL            string
}

func NewGetInfoUseCase(youtubeRepository youtube.IRepository) *GetInfoUseCase {
	return &GetInfoUseCase{youtubeRepository: youtubeRepository}
}

func (uc *GetInfoUseCase) Execute(ctx context.Context, in *GetInfoUseCaseInput) (*GetInfoUseCaseOutput, error) {
	channel, err := uc.youtubeRepository.FindByID(ctx, in.ChannelID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	videos := make([]*Video, 0, len(channel.Videos))
	for _, video := range channel.Videos {
		videos = append(videos, &Video{
			Number:         video.GetNumber(),
			Title:          video.Title,
			MinuteDuration: int(math.Round(video.Duration.Minutes())),
			URL:            video.GetURL(),
		})
	}

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Number < videos[j].Number
	})

	return &GetInfoUseCaseOutput{Videos: videos}, nil
}
