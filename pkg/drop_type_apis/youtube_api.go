package drop_type_apis

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
)

type YoutubeAPI struct {
	ContentTitle       string
	ContentDescription string
	ContentPicturePath string
	Client             *youtube.Service
}

func (y *YoutubeAPI) Search(search string) []ApiSearchResponse {
	call := y.Client.Search.List([]string{"snippet"}).Q(search).MaxResults(20).Type("video")

	response, err := call.Do()

	if err != nil {
		fmt.Printf("Unable to retrieve search results with search query %s : %v\n", search, err)
		return nil
	}

	if len(response.Items) == 0 {
		fmt.Printf("No results found for search query %s\n", search)
		return nil
	}

	var results []ApiSearchResponse
	for _, item := range response.Items {
		results = append(results, ApiSearchResponse{
			Search:      search,
			PicturePath: item.Snippet.Thumbnails.Default.Url,
			Title:       item.Snippet.Title,
			Subtitle:    item.Snippet.Description,
		})
	}

	return results
}

func (y *YoutubeAPI) Init() {
	apiKey := os.Getenv("YOUTUBE_API_KEY")

	if apiKey == "" {
		fmt.Printf("YouTube API key not found in environment variable YOUTUBE_API_KEY\n")
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("Unable to create YouTube service: %v\n", err)
	}
	y.Client = service
}

var (
	_ DropTypeAPI = &YoutubeAPI{}
)

const credentialFile = "config/client_secrets.json"
const tokenFile = "token.json"

func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
