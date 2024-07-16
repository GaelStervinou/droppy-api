package drop_type_apis

func Factory(dropType string) DropTypeAPI {
	switch dropType {
	case YoutubeType:
		return &YoutubeAPI{}
	case SpotifyType:
		return &SpotifyAPI{}
	default:
		return nil
	}
}
