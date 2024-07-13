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
	GetFollowed() UserModel
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
	IsActiveFollowing(followerID, followedID uint) (bool, error)
	IsPendingFollowing(followerID, followedID uint) (bool, error)
	CountFollowers(userID uint) int
	CountFollowed(userID uint) int
	GetUserFollowedBy(followerID uint, followedID uint) (FollowModel, error)
}

type FollowService interface {
	GetUserFollowing(userID uint, requesterID uint) ([]FollowModel, error)
	GetUserFollowers(userID uint, requesterID uint) ([]FollowModel, error)
}

type FollowCreationParam struct {
	UserToFollowID uint `json:"userId"`
}
