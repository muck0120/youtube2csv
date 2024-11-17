package main

import (
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("\n👇 Please access the displayed URL to obtain the authentication code and enter it.\n\n%v\n\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		slog.Error("unable to read authorization code", slog.String("error", err.Error()))
		os.Exit(1)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		slog.Error("unable to retrieve token from web", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return token
}

func channelsListByUsername(service *youtube.Service, part string, forUsername string) {
	call := service.Channels.List([]string{part})
	call = call.ForUsername(forUsername)

	response, err := call.Do()
	if err != nil {
		slog.Error("failed to call YouTube API", slog.String("error", err.Error()))
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s' and it has %d views.", response.Items[0].Id, response.Items[0].Snippet.Title, response.Items[0].Statistics.ViewCount))
}

func main() {
	ctx := context.Background()

	secret, err := os.ReadFile(os.Getenv("WORKDIR") + "/input/client_secret.json")
	if err != nil {
		slog.Error("unable to read client secret file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	config, err := google.ConfigFromJSON(secret, youtube.YoutubeReadonlyScope)
	if err != nil {
		slog.Error("unable to parse client secret file to config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	token := getTokenFromWeb(ctx, config)

	service, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		slog.Error("failed to create YouTube service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	channelsListByUsername(service, "snippet,contentDetails,statistics", "GoogleDevelopers")
}
