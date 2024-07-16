package drop_type_apis

import (
	"fmt"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"os"
)

var (
	_ DropTypeAPI = &SpotifyAPI{}
)

type SpotifyAPI struct {
	Client *spotify.Client
}

func (s *SpotifyAPI) Search(search string) []ApiSearchResponse {
	result, err := s.Client.Search(context.Background(), search, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatalf("Failed to search for track: %v", err)
		return nil
	}

	if len(result.Tracks.Tracks) == 0 {
		log.Printf("No results found for search query %s\n", search)
		return nil
	}
	fmt.Println(result.Tracks.Tracks[0].Name)

	var results []ApiSearchResponse
	for _, item := range result.Tracks.Tracks {
		results = append(results, ApiSearchResponse{
			Search:      search,
			PicturePath: item.Album.Images[0].URL,
			Title:       item.Name,
			Subtitle:    item.Artists[0].Name,
			Content:     string(item.URI),
		})
	}
	return results
}

func (s *SpotifyAPI) Init() {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	// Set up the OAuth2 config
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	ctx := context.Background()
	// Get a token and create a Spotify client
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	s.Client = spotify.New(httpClient)
}
