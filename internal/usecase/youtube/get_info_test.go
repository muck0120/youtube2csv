package youtube_test

import (
	"context"
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	domain "github.com/muck0120/youtube2csv/internal/domain/youtube"
	"github.com/muck0120/youtube2csv/internal/domain/youtube/fixture"
	usecase "github.com/muck0120/youtube2csv/internal/usecase/youtube"
)

func TestGetInfoUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arg       *usecase.GetInfoUseCaseInput
		setup     func(*domain.MockIRepository)
		want      *usecase.GetInfoUseCaseOutput
		wantError bool
	}{
		{
			name: "期待したレスポンスが取得できる",
			arg:  &usecase.GetInfoUseCaseInput{ChannelID: "ABCDEFG123456"},
			setup: func(repo *domain.MockIRepository) {
				repo.EXPECT().FindByID(gomock.Any(), "ABCDEFG123456").Return(fixture.Chnanel(nil), nil)
			},
			want: func() *usecase.GetInfoUseCaseOutput {
				channel := fixture.Chnanel(nil)

				videos := make([]*usecase.Video, 0, len(channel.Videos))
				for _, video := range channel.Videos {
					videos = append(videos, &usecase.Video{
						Number:         video.GetNumber(),
						Title:          video.Title,
						MinuteDuration: int(math.Round(video.Duration.Minutes())),
						URL:            video.GetURL(),
					})
				}

				sort.Slice(videos, func(i, j int) bool {
					return videos[i].Number < videos[j].Number
				})

				return &usecase.GetInfoUseCaseOutput{Videos: videos}
			}(),
			wantError: false,
		},
		{
			name: "Channel が見つからない場合、エラーが返る",
			arg:  &usecase.GetInfoUseCaseInput{ChannelID: "ABCDEFG123456"},
			setup: func(repo *domain.MockIRepository) {
				repo.EXPECT().FindByID(gomock.Any(), "ABCDEFG123456").Return(nil, assert.AnError)
			},
			want:      nil,
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := domain.NewMockIRepository(ctrl)

			test.setup(repo)

			got, err := usecase.NewGetInfoUseCase(repo).Execute(context.Background(), test.arg)

			assert.Equal(t, test.wantError, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}
