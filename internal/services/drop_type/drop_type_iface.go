package drop_type

type DropType interface {
	Init() error
	IsValidContent(string) bool
}
