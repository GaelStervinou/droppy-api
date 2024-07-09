package model

type LikeModel interface {
	GetID() uint
	GetDropID() uint
	GetDrop() DropModel
	GetUserID() uint
	GetUser() UserModel
}

type LikeRepository interface {
	CreateLike(dropId uint, userId uint) (LikeModel, error)
	DeleteLike(dropId uint, userId uint) error
	GetDropTotalLikes(dropId uint) (int, error)
	LikeExists(dropId uint, userId uint) (bool, error)
}

type LikeService interface {
	LikeDrop(userID uint, args LikeParam) (LikeModel, error)
	UnlikeDrop(userID uint, args LikeParam) error
}

type LikeParam struct {
	DropId uint `json:"dropId"`
}
