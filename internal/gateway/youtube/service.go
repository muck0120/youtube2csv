package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/muck0120/youtube2csv/internal/pkg/errors"
)

type IService interface {
	FindVideosByChannelID(id string) ([]*youtube.Video, error)
}

type Service struct {
	client *youtube.Service
}

func NewService(ctx context.Context) (*Service, error) {
	return sync.OnceValues(func() (*Service, error) {
		secret, err := os.ReadFile(os.Getenv("WORKDIR") + "/input/client_secret.json")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		config, err := google.ConfigFromJSON(secret, youtube.YoutubeReadonlyScope)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		token, err := getToken(ctx, config)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		client, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return &Service{client: client}, nil
	})()
}

func getToken(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	token, err := getTokenFromCache()
	if err == nil { // キャッシュファイルがあればそれを返す
		return token, nil
	}

	token, err = getTokenFromWeb(ctx, config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = saveToken(token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return token, nil
}

func getTokenFromCache() (*oauth2.Token, error) {
	file, err := os.Open(os.Getenv("WORKDIR") + "/input/token.json")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()

	token := &oauth2.Token{}
	if err = json.NewDecoder(file).Decode(token); err != nil {
		return nil, errors.WithStack(err)
	}

	return token, nil
}

func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("\n👇 Please access the displayed URL to obtain the authentication code and enter it.\n\n%v\n\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, errors.WithStack(err)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return token, nil
}

func saveToken(token *oauth2.Token) error {
	file, err := os.OpenFile((os.Getenv("WORKDIR") + "/input/token.json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(token); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (sv *Service) FindVideosByChannelID(channelID string) ([]*youtube.Video, error) {
	channelResponse, err := sv.client.Channels.List([]string{"contentDetails"}).Id(channelID).Do()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(channelResponse.Items) == 0 {
		return []*youtube.Video{}, nil
	}

	uploadsPlaylistID := channelResponse.Items[0].ContentDetails.RelatedPlaylists.Uploads

	videoIDs := make([]string, 0)
	nextPageToken := ""

	for {
		var maxResults int64 = 1000
		call := sv.client.PlaylistItems.List([]string{"contentDetails"}).PlaylistId(uploadsPlaylistID).MaxResults(maxResults)

		if nextPageToken != "" {
			call = call.PageToken(nextPageToken)
		}

		playlistResponse, err := call.Do()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, video := range playlistResponse.Items {
			videoIDs = append(videoIDs, video.ContentDetails.VideoId)
		}

		if playlistResponse.NextPageToken == "" {
			break
		}

		nextPageToken = playlistResponse.NextPageToken
	}

	res, per := make([]*youtube.Video, 0, len(videoIDs)), 50
	for chunk := range slices.Chunk(videoIDs, per) {
		videoResponse, err := sv.client.Videos.List([]string{"contentDetails", "snippet", "status"}).Id(chunk...).Do()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		res = append(res, videoResponse.Items...)
	}

	return res, nil
}
