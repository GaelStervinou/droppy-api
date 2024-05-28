package drop_type

func NewDropTypeFactory() *DropTypeFactory {
	return &DropTypeFactory{}
}

type DropTypeFactory struct{}

func (d *DropTypeFactory) CreateDropType(dropType string) DropType {
	switch dropType {
	case "youtube":
		return &Youtube{}
	default:
		return nil
	}
}
