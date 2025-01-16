package main

import (
	"context"
	"flag"
	"log/slog"
	"time"

	controller "github.com/muck0120/youtube2csv/internal/controller/youtube"
	"github.com/muck0120/youtube2csv/internal/gateway/youtube"
	"github.com/muck0120/youtube2csv/internal/pkg/errors"
	usecase "github.com/muck0120/youtube2csv/internal/usecase/youtube"
)

func main() {
	ctx := context.Background()

	slog.Info("start youtube info to csv")

	defer func() {
		slog.Info("end youtube info to csv")
	}()

	var (
		secretFile = flag.String("secret", "./client_secret.json", "file path to client_secret.json")
		tokenFile  = flag.String("token", "./token.json", "where to store access token (or where to read from)")
		channelID  = flag.String("channel-id", "", "target channel ID")
		outputFile = flag.String("out", "", "output file path")
	)

	flag.Parse()

	if *channelID == "" {
		slog.Error("Channel ID is required")

		return
	}

	if *outputFile == "" {
		*outputFile = ("./output/" + *channelID + "_" + time.Now().Format("20060102150405") + ".csv")
	}

	service, err := youtube.NewService(ctx, *secretFile, *tokenFile)
	if err != nil {
		slog.Error("failed to create YouTube service", slog.String("error", err.Error()), errors.LogStackTrace(err))

		return
	}

	rp := youtube.NewRepository(service)
	uc := usecase.NewGetInfoUseCase(rp)
	in := &controller.GetInfoControllerInput{ChannelID: *channelID, FilePath: *outputFile}

	if err := controller.NewGetInfoController(uc).Run(ctx, in); err != nil {
		slog.Error("failed to run GetInfoController", slog.String("error", err.Error()), errors.LogStackTrace(err))

		return
	}
}
