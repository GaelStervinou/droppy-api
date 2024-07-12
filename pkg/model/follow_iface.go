package model

type FollowStatus interface {
	ToInt() uint
}

type FollowModel interface {
	GetID() uint
	GetFollowerID() uint
	GetFollowedID() uint
	GetStatus() uint
	GetCreatedAt() uint
	GetFollower() UserModel
}

type FollowRepository interface {
	Create(followerID, followedID uint, isActive bool) (FollowModel, error)
	AcceptRequest(followId uint) error
	RejectRequest(followId uint) error
	Delete(followId uint) error
	GetPendingRequests(userID uint) ([]FollowModel, error)
	GetFollowers(userID uint) ([]FollowModel, error)
	GetFollowing(userID uint) ([]FollowModel, error)
	AreAlreadyFollowing(followerID, followedID uint) (bool, error)
	IsFollowing(followerID, followedID uint) (bool, error)
	CountFollowers(userID uint) int
	CountFollowed(userID uint) int
}

type FollowCreationParam struct {
	UserToFollowID uint `json:"userId"`
}
