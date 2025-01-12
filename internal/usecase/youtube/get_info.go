package youtube

import (
	"context"

	"muck0120/youtube2csv/internal/domain/youtube"
	"muck0120/youtube2csv/pkg/errors"
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
	Channel *youtube.Channel
}

func NewGetInfoUseCase(youtubeRepository youtube.IRepository) *GetInfoUseCase {
	return &GetInfoUseCase{youtubeRepository: youtubeRepository}
}

func (uc *GetInfoUseCase) Execute(ctx context.Context, in *GetInfoUseCaseInput) (*GetInfoUseCaseOutput, error) {
	channel, err := uc.youtubeRepository.FindByID(ctx, in.ChannelID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &GetInfoUseCaseOutput{Channel: channel}, nil
}
