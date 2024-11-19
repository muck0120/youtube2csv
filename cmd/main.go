package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func getToken(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	token, err := getTokenFromCache()
	if err != nil {
		token = getTokenFromWeb(ctx, config)
		saveToken(token)
	}

	return token
}

func getTokenFromCache() (*oauth2.Token, error) {
	file, err := os.Open(os.Getenv("WORKDIR") + "/input/token.json")
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	token := &oauth2.Token{}
	if err = json.NewDecoder(file).Decode(token); err != nil {
		return nil, fmt.Errorf("unable to decode token: %w", err)
	}

	return token, nil
}

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

func saveToken(token *oauth2.Token) {
	file, err := os.OpenFile((os.Getenv("WORKDIR") + "/input/token.json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		slog.Error("unable to cache oauth token", slog.String("error", err.Error()))
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(token); err != nil {
		slog.Error("unable to cache oauth token", slog.String("error", err.Error()))
	}
}

func channelsList(service *youtube.Service) {
	var count int64 = 300
	call := service.Search.List([]string{"snippet"}).ChannelId(os.Getenv("CHANNEL_ID")).MaxResults(count)

	res, err := call.Do()
	if err != nil {
		slog.Error("failed to call YouTube API", slog.String("error", err.Error()))
		os.Exit(1)
	}

	for _, item := range res.Items {
		if item.Id.Kind == "youtube#video" {
			fmt.Printf("ID: %s\n", item.Id.VideoId)
			call2 := service.Videos.List([]string{"snippet"}).Id(item.Id.VideoId)
			res2, _ := call2.Do()

			for _, video := range res2.Items {
				fmt.Printf("Title: %s\n", video.Snippet.Title)
			}
		}
	}
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

	token := getToken(ctx, config)

	service, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		slog.Error("failed to create YouTube service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	channelsList(service)
}
