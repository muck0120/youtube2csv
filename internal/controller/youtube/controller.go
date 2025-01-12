package youtube

import (
	"context"
	"fmt"

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

	c.print(ctx, out)

	return nil
}

func (c *GetInfoController) print(_ context.Context, out *youtube.GetInfoUseCaseOutput) {
	for _, video := range out.Channel.Videos {
		fmt.Printf("Title: %s\n", video.Title)
		fmt.Printf("Duration: %s\n", video.Duration)
		fmt.Printf("Status: %s\n", video.Status)
	}
}
