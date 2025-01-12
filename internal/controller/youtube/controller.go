package youtube

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"muck0120/youtube2csv/internal/usecase/youtube"
	"muck0120/youtube2csv/pkg/errors"
)

type IGetInfoController interface {
	Run(ctx context.Context, in *GetInfoControllerInput) error
}

type GetInfoController struct {
	GetInfoUseCase youtube.IGetInfoUseCase
}

type GetInfoControllerInput struct {
	ChannelID string
}

func NewGetInfoController(getInfoUseCase youtube.IGetInfoUseCase) *GetInfoController {
	return &GetInfoController{GetInfoUseCase: getInfoUseCase}
}

func (c *GetInfoController) Run(ctx context.Context, in *GetInfoControllerInput) error {
	out, err := c.GetInfoUseCase.Execute(ctx, &youtube.GetInfoUseCaseInput{ChannelID: in.ChannelID})
	if err != nil {
		return errors.WithStack(err)
	}

	if err := c.csv(ctx, out); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *GetInfoController) csv(_ context.Context, out *youtube.GetInfoUseCaseOutput) error {
	file, err := os.Create(os.Getenv("WORKDIR") + "/output/" + time.Now().Format("20060102150405") + ".csv")
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
