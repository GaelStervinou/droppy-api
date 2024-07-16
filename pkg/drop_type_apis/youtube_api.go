package drop_type_apis

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"os"
)

var (
	_ DropTypeAPI = &YoutubeAPI{}
)

type YoutubeAPI struct {
	Client *youtube.Service
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
			Subtitle:    item.Snippet.ChannelTitle,
			Content:     generateUrl(item.Id.VideoId),
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

func generateUrl(videoId string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)
}
