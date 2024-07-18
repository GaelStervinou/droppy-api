package postgres

import (
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	FollowerID uint
	Follower   User `gorm:"foreignKey:FollowerID"`
	FollowedID uint
	Followed   User `gorm:"foreignKey:FollowedID"`
	Status     uint
}

func (f *Follow) GetID() uint {
	return f.ID
}

func (f *Follow) GetFollowerID() uint {
	return f.FollowerID
}

func (f *Follow) GetFollowedID() uint {
	return f.FollowedID
}

func (f *Follow) GetStatus() uint {
	return f.Status
}

func (f *Follow) GetCreatedAt() uint {
	return uint(f.CreatedAt.Unix())
}

func (f *Follow) GetFollower() model.UserModel {
	return &f.Follower
}

func (f *Follow) GetFollowed() model.UserModel {
	return &f.Followed
}

type FollowPendingStatus struct {
}

func (f *FollowPendingStatus) ToInt() uint {
	return 0
}

type FollowAcceptedStatus struct {
}

func (f *FollowAcceptedStatus) ToInt() uint {
	return 1
}

var _ model.FollowModel = (*Follow)(nil)

type repoFollowPrivate struct {
	db *gorm.DB
}

func NewFollowRepo(db *gorm.DB) model.FollowRepository {
	return &repoFollowPrivate{db: db}
}

func (r *repoFollowPrivate) Create(followerID, followedID uint, isPublic bool) (model.FollowModel, error) {
	var status uint
	if isPublic {
		status = new(FollowAcceptedStatus).ToInt()
	} else {
		status = new(FollowPendingStatus).ToInt()
	}
	follow := &Follow{
		FollowerID: followerID,
		FollowedID: followedID,
		Status:     status,
	}
	result := r.db.Create(follow)
	if result.Error != nil {
		return nil, result.Error
	}
	return follow, nil
}

func (r *repoFollowPrivate) AcceptRequest(followId uint) error {
	result := r.db.Model(&Follow{}).Where("id = ?", followId).Update("status", new(FollowAcceptedStatus).ToInt())
	return result.Error
}

func (r *repoFollowPrivate) RejectRequest(followId uint) error {
	return r.Delete(followId)
}

func (r *repoFollowPrivate) Delete(followId uint) error {
	result := r.db.Delete(&Follow{}, followId)
	return result.Error
}

func (r *repoFollowPrivate) GetPendingRequests(userID uint) ([]model.FollowModel, error) {
	var follows []Follow
	result := r.db.
		Preload("Follower").
		Where("followed_id = ? AND status = ?", userID, new(FollowPendingStatus).ToInt()).Find(&follows)
	if result.Error != nil {
		return nil, result.Error
	}
	var models []model.FollowModel
	for _, follow := range follows {
		models = append(models, &follow)
	}
	return models, nil
}

func (r *repoFollowPrivate) GetFollowers(userID uint) ([]model.FollowModel, error) {
	var follows []Follow
	result := r.db.
		Preload("Followed").
		Preload("Follower").
		Joins("JOIN users AS follower ON follower.id = follows.follower_id AND follower.status = ? AND follower.deleted_at IS NULL", 1).
		Where("followed_id = ? AND follows.status = ?", userID, new(FollowAcceptedStatus).ToInt()).Find(&follows)
	if result.Error != nil {
		return nil, result.Error
	}
	var models []model.FollowModel
	for _, follow := range follows {
		models = append(models, &follow)
	}
	return models, nil
}

func (r *repoFollowPrivate) GetFollowing(userID uint) ([]model.FollowModel, error) {
	var follows []Follow
	result := r.db.
		Preload("Followed").
		Preload("Follower").
		Joins("JOIN users AS followed ON followed.id = follows.followed_id AND followed.status = ? AND followed.deleted_at IS NULL", 1).
		Where("follower_id = ? AND follows.status = ?", userID, new(FollowAcceptedStatus).ToInt()).Find(&follows)
	if result.Error != nil {
		return nil, result.Error
	}
	var models []model.FollowModel
	for _, follow := range follows {
		models = append(models, &follow)
	}
	return models, nil
}

func (r *repoFollowPrivate) AreAlreadyFollowing(followerID, followedID uint) (bool, error) {
	var follow Follow
	result := r.db.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (r *repoFollowPrivate) IsActiveFollowing(followerID, followedID uint) (bool, error) {
	var follow Follow
	result := r.db.Where("follower_id = ? AND followed_id = ? AND status = ?", followerID, followedID, new(FollowAcceptedStatus).ToInt()).Find(&follow)
	if result.Error != nil {
		return false, result.Error
	}
	if follow.CreatedAt.IsZero() {
		return false, nil
	}
	return true, nil
}

func (r *repoFollowPrivate) IsPendingFollowing(followerID, followedID uint) (bool, error) {
	var follow Follow
	result := r.db.Where("follower_id = ? AND followed_id = ? AND status = ?", followerID, followedID, new(FollowPendingStatus).ToInt()).Find(&follow)
	if result.Error != nil {
		return false, result.Error
	}
	if follow.CreatedAt.IsZero() {
		return false, nil
	}
	return true, nil
}

func (r *repoFollowPrivate) CountFollowers(userID uint) int {
	var count int64
	result := r.db.Table("follows").
		Joins("JOIN users ON users.id = follows.follower_id").
		Where("follows.deleted_at IS NULL AND follows.status = ? AND follows.followed_id = ? AND users.status = ? AND users.deleted_at IS NULL", 1, userID, 1).
		Count(&count).Count(&count)
	if result.Error != nil {
		return 0
	}
	return int(count)
}

func (r *repoFollowPrivate) CountFollowed(userID uint) int {
	var count int64
	result := r.db.Model(&Follow{}).Where("follower_id = ? AND status = ?", userID, new(FollowAcceptedStatus).ToInt()).Count(&count)
	if result.Error != nil {
		return 0
	}
	return int(count)
}

func (r *repoFollowPrivate) GetUserFollowedBy(followerID uint, followedID uint) (model.FollowModel, error) {
	var follow Follow
	result := r.db.
		Preload("Follower").
		Preload("Followed").
		Where("follower_id = ? AND followed_id = ?", followerID, followedID).Find(&follow)
	if result.Error != nil {
		return nil, result.Error
	}

	if follow.CreatedAt.IsZero() {
		return nil, nil
	}

	return &follow, nil
}

func (r *repoFollowPrivate) GetFollowByID(followID uint) (model.FollowModel, error) {
	var follow Follow
	result := r.db.
		Preload("Follower").
		Preload("Followed").
		Where("id = ?", followID).First(&follow)
	if result.Error != nil {
		return nil, result.Error
	}

	if follow.CreatedAt.IsZero() {
		return nil, nil
	}

	return &follow, nil
}

func (r *repoFollowPrivate) GetPendingFollowByID(followID uint) (model.FollowModel, error) {
	var follow Follow
	result := r.db.
		Preload("Follower").
		Preload("Followed").
		Where("id = ? AND status = ?", followID, new(FollowPendingStatus).ToInt()).First(&follow)
	if result.Error != nil {
		return nil, result.Error
	}

	if follow.CreatedAt.IsZero() {
		return nil, nil
	}

	return &follow, nil
}
