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
	GetLat() float64
	GetLng() float64
	GetPicturePath() string
	GetCreatedAt() int
	GetCreatedBy() UserModel
}

type DropRepository interface {
	Create(dropNotificationId uint, contentType string, content string, description string, createdById uint, status uint, isPinned bool, picturePath string, lat float64, lng float64) (DropModel, error)
	Delete(dropId uint) error
	GetUserDrops(userId uint) ([]DropModel, error)
	GetDropByDropNotificationAndUser(dropNotificationId uint, userId uint) (DropModel, error)
	GetDropsByUserIdsAndDropNotificationId(userIds []uint, dropNotifId uint) ([]DropModel, error)
	HasUserDropped(dropNotificationId uint, userId uint) (bool, error)
	GetDropById(dropId uint) (DropModel, error)
}

type DropService interface {
	CanCreateDrop(userId uint) (bool, error)
	IsValidDropCreation(args DropCreationParam) (bool, error)
	CreateDrop(userId uint, args DropCreationParam) (DropModel, error)
	GetUserFeed(userId uint) ([]DropModel, error)
	GetDropsByUserId(userId uint, currentUser UserModel) ([]DropModel, error)
	HasUserDroppedToday(userId uint) (bool, error)
}

type DropCreationParam struct {
	Content     string  `json:"content"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Picture     string  `json:"picture"`
}

type FilledDropCreation struct {
	Type               string  `json:"type"`
	Content            string  `json:"content"`
	Description        string  `json:"description"`
	DropNotificationId uint    `json:"dropNotificationId"`
	PicturePath        string  `json:"picturePath"`
	Lat                float64 `json:"lat"`
	Lng                float64 `json:"lng"`
}
