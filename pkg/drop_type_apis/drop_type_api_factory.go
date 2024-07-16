package drop_type_apis

func Factory(dropType string) DropTypeAPI {
	switch dropType {
	case YoutubeType:
		return &YoutubeAPI{}
	case SpotifyType:
		return &SpotifyAPI{}
	case FilmType:
		return &FilmsAPI{}
	case TwitchType:
		return &TwitchTypeApi{}
	default:
		return nil
	}
}
