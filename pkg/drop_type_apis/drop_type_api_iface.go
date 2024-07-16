package drop_type_apis

type DropTypeAPI interface {
	Search(search string) []ApiSearchResponse
	Init()
}

type ApiSearch interface {
	GetSearch() string
	GetContentTitle() string
	GetContentDescription() string
	GetContentPicturePath() string
	GetContent() string
}

type ApiSearchResponse struct {
	Search      string
	PicturePath string
	Title       string
	Subtitle    string
	Content     string
}

func (a ApiSearchResponse) GetSearch() string {
	return a.Search
}

func (a ApiSearchResponse) GetContentTitle() string {
	return a.Title
}

func (a ApiSearchResponse) GetContentDescription() string {
	return a.Subtitle
}

func (a ApiSearchResponse) GetContentPicturePath() string {
	return a.PicturePath
}

func (a ApiSearchResponse) GetContent() string {
	return a.Content
}
