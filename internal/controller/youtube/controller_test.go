package youtube_test

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	controller "github.com/muck0120/youtube2csv/internal/controller/youtube"
	usecase "github.com/muck0120/youtube2csv/internal/usecase/youtube"
)

func TestGetInfoController_Run(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		arg       *controller.GetInfoControllerInput
		setup     func(*usecase.MockIGetInfoUseCase)
		wantError bool
		more      func(*testing.T)
	}{
		{
			name: "期待した内容が CSV ファイルに書き込まれる",
			arg: &controller.GetInfoControllerInput{
				ChannelID: "ABCDEFG123456",
				FilePath:  tmpDir + "test1.csv",
			},
			setup: func(repo *usecase.MockIGetInfoUseCase) {
				input := &usecase.GetInfoUseCaseInput{ChannelID: "ABCDEFG123456"}

				repo.EXPECT().Execute(gomock.Any(), input).Return(&usecase.GetInfoUseCaseOutput{
					Videos: []*usecase.Video{
						{
							Number:         1,
							Title:          "title1",
							MinuteDuration: 1,
							URL:            "https://example.com/1",
						},
						{
							Number:         2,
							Title:          "title2",
							MinuteDuration: 2,
							URL:            "https://example.com/2",
						},
						{
							Number:         3,
							Title:          "title3",
							MinuteDuration: 3,
							URL:            "https://example.com/3",
						},
					},
				}, nil)
			},
			wantError: false,
			more: func(t *testing.T) {
				t.Helper()

				// expected
				file1, err := os.Open("./testdata/test1.csv")
				if err != nil {
					t.Error(err)

					return
				}
				defer file1.Close()

				reader1 := csv.NewReader(file1)
				csv1, err := reader1.ReadAll()
				if err != nil {
					t.Error(err)

					return
				}

				// actual
				file2, err := os.Open((tmpDir + "test1.csv"))
				if err != nil {
					t.Error(err)

					return
				}
				defer file2.Close()

				reader2 := csv.NewReader(file2)
				csv2, err := reader2.ReadAll()
				if err != nil {
					t.Error(err)

					return
				}

				assert.Equal(t, csv1, csv2)
			},
		},
		{
			name: "ユースケースでエラーが発生した場合、エラーが返る",
			arg: &controller.GetInfoControllerInput{
				ChannelID: "ABCDEFG123456",
				FilePath:  tmpDir + "test2.csv",
			},
			setup: func(repo *usecase.MockIGetInfoUseCase) {
				input := &usecase.GetInfoUseCaseInput{ChannelID: "ABCDEFG123456"}
				repo.EXPECT().Execute(gomock.Any(), input).Return(nil, assert.AnError)
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := usecase.NewMockIGetInfoUseCase(ctrl)

			test.setup(repo)

			err := controller.NewGetInfoController(repo).Run(t.Context(), test.arg)
			assert.Equal(t, test.wantError, err != nil)

			if test.more != nil {
				test.more(t)
			}
		})
	}
}
