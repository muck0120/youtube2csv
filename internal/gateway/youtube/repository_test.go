package youtube_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	domain "github.com/muck0120/youtube2csv/internal/domain/youtube"
	domainFixture "github.com/muck0120/youtube2csv/internal/domain/youtube/fixture"
	gateway "github.com/muck0120/youtube2csv/internal/gateway/youtube"
	gatewayFixture "github.com/muck0120/youtube2csv/internal/gateway/youtube/fixture"
	"github.com/muck0120/youtube2csv/internal/pkg/time"
)

func TestRepository_FindByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arg       string
		setup     func(*gateway.MockIService)
		want      *domain.Channel
		wantError bool
	}{
		{
			name: "期待したレスポンスが取得できる",
			arg:  "ABCDEFG123456",
			setup: func(srv *gateway.MockIService) {
				srv.EXPECT().FindVideosByChannelID("ABCDEFG123456").Return(gatewayFixture.Videos(nil), nil)
			},
			want: func() *domain.Channel {
				return domainFixture.Chnanel(func(channel *domain.Channel) {
					channel.ID = "ABCDEFG123456"

					videos := gatewayFixture.Videos(nil)

					ret := make([]*domain.Video, 0, len(videos))
					for _, video := range videos {
						duration, _ := time.ParseISO8601Duration(video.ContentDetails.Duration)

						ret = append(ret, &domain.Video{
							ID:       video.Id,
							Title:    video.Snippet.Title,
							Duration: duration,
							Status:   domain.VideoStatus(video.Status.PrivacyStatus),
						})
					}

					channel.Videos = ret
				})
			}(),
			wantError: false,
		},
		{
			name: "YouTube API からのレスポンスがエラーの場合、エラーが返る",
			arg:  "ABCDEFG123456",
			setup: func(srv *gateway.MockIService) {
				srv.EXPECT().FindVideosByChannelID("ABCDEFG123456").Return(nil, assert.AnError)
			},
			want:      nil,
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			srv := gateway.NewMockIService(ctrl)

			test.setup(srv)

			got, err := gateway.NewRepository(srv).FindByID(context.Background(), test.arg)

			assert.Equal(t, test.wantError, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}
