package drop_type

import (
	"context"
	"google.golang.org/api/youtube/v3"
)

type Youtube struct{}

func (y *Youtube) Init() error {
	return nil
}

func (y *Youtube) IsValidContent(uri string) bool {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx)

	if err != nil {
		return false
	}

	_, err = youtubeService.Videos.List([]string{"id"}).Id(uri).Do()

	if err != nil {
		return false
	}

	return true
}
