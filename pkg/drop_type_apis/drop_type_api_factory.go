package drop_type_apis

func Factory(dropType string) DropTypeAPI {
	switch dropType {
	case "youtube":
		return &YoutubeAPI{}
	default:
		return nil
	}
}
