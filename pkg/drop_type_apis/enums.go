package drop_type_apis

import "slices"

var YoutubeType = "youtube"
var SpotifyType = "spotify"
var FilmType = "films"
var TwitchType = "twitch"

var ValidDropTypes = []string{YoutubeType, SpotifyType, FilmType, TwitchType}

func GetValidDropTypes() []string {
	return ValidDropTypes
}

func IsValidDropType(dropType string) bool {
	return slices.Contains(ValidDropTypes, dropType)
}
