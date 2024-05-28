package model

type DropStatus interface {
	ToInt() int
}

type DropModel interface {
	GetID() uint
	GetType() string
	GetContent() string
	GetDescription() string
	GetCreatedById() uint
	GetStatus() uint
	GetDeletedById() uint
	GetIsPinned() bool
	GetDropNotificationID() uint
}

type DropRepository interface {
	Create(dropNotificationId uint, contentType string, content string, description string, createdById uint, status uint, isPinned bool) (DropModel, error)
	Delete(dropId uint) error
	GetUserDrops(userId uint) ([]DropModel, error)
	GetDropByDropNotificationAndUser(dropNotificationId uint, userId uint) (DropModel, error)
}

type DropService interface {
	CanCreateDrop(current, userId uint) (bool, error)
	IsValidDropCreation(args DropCreationParam) (bool, error)
	CreateDrop(userId uint, args DropCreationParam) (DropModel, error)
}

type DropCreationParam struct {
	Type               string `json:"type"`
	Content            string `json:"content"`
	Description        string `json:"description"`
	DropNotificationId uint   `json:"dropNotificationId"`
}
