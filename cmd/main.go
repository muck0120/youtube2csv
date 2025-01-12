package main

import (
	"context"
	"log/slog"
	"os"

	controller "muck0120/youtube2csv/internal/controller/youtube"
	"muck0120/youtube2csv/internal/gateway/youtube"
	usecase "muck0120/youtube2csv/internal/usecase/youtube"
	"muck0120/youtube2csv/pkg/errors"
)

func main() {
	ctx := context.Background()

	slog.Info("start youtube info to csv")

	defer func() {
		slog.Info("end youtube info to csv")
	}()

	service, err := youtube.NewService(ctx)
	if err != nil {
		slog.Error("failed to create YouTube service", slog.String("error", err.Error()))

		return
	}

	rp := youtube.NewRepository(service)
	uc := usecase.NewGetInfoUseCase(rp)
	in := &controller.GetInfoControllerInput{ChannelID: os.Getenv("CHANNEL_ID")}

	if err := controller.NewGetInfoController(uc).Run(ctx, in); err != nil {
		slog.Error("failed to run GetInfoController", slog.String("error", err.Error()), errors.LogStackTrace(err))

		return
	}
}
