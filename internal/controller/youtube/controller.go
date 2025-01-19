//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=$GOPACKAGE -destination=./controller_mock.go

package youtube

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/muck0120/youtube2csv/internal/pkg/errors"
	"github.com/muck0120/youtube2csv/internal/usecase/youtube"
)

type IGetInfoController interface {
	Run(ctx context.Context, in *GetInfoControllerInput) error
}

type GetInfoController struct {
	GetInfoUseCase youtube.IGetInfoUseCase
}

type GetInfoControllerInput struct {
	ChannelID string
	FilePath  string
}

func NewGetInfoController(getInfoUseCase youtube.IGetInfoUseCase) *GetInfoController {
	return &GetInfoController{GetInfoUseCase: getInfoUseCase}
}

func (c *GetInfoController) Run(ctx context.Context, input *GetInfoControllerInput) error {
	out, err := c.GetInfoUseCase.Execute(ctx, &youtube.GetInfoUseCaseInput{ChannelID: input.ChannelID})
	if err != nil {
		return errors.WithStack(err)
	}

	if err := c.csv(ctx, input.FilePath, out); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *GetInfoController) csv(_ context.Context, path string, out *youtube.GetInfoUseCaseOutput) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダーを書き込む
	header := []string{"no", "title", "duration", "url"}
	if err := writer.Write(header); err != nil {
		return errors.WithStack(err)
	}

	for _, video := range out.Videos {
		row := []string{
			strconv.Itoa(video.Number),         // 動画番号
			video.Title,                        // 動画タイトル
			strconv.Itoa(video.MinuteDuration), // 動画時間
			video.URL,                          // 動画URL
		}

		if err := writer.Write(row); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
