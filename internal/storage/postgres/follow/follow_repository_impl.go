package follow

import (
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	FollowerID uint
	FollowedID uint
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

type repoPrivate struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) model.FollowRepository {
	return &repoPrivate{db: db}
}

func (r *repoPrivate) Create(followerID, followedID uint, isActive bool) (model.FollowModel, error) {
	var status uint
	if isActive {
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

func (r *repoPrivate) AcceptRequest(followId uint) error {
	result := r.db.Model(&Follow{}).Where("id = ?", followId).Update("status", new(FollowAcceptedStatus).ToInt())
	return result.Error
}

func (r *repoPrivate) RejectRequest(followId uint) error {
	return r.Delete(followId)
}

func (r *repoPrivate) Delete(followId uint) error {
	result := r.db.Delete(&Follow{}, followId)
	return result.Error
}

func (r *repoPrivate) GetPendingRequests(userID uint) ([]model.FollowModel, error) {
	var follows []Follow
	result := r.db.Where("followed_id = ? AND status = ?", userID, new(FollowPendingStatus).ToInt()).Find(&follows)
	if result.Error != nil {
		return nil, result.Error
	}
	var models []model.FollowModel
	for _, follow := range follows {
		models = append(models, &follow)
	}
	return models, nil
}

func (r *repoPrivate) GetFollowers(userID uint) ([]model.FollowModel, error) {
	var follows []Follow
	result := r.db.Where("followed_id = ? AND status = ?", userID, new(FollowAcceptedStatus).ToInt()).Find(&follows)
	if result.Error != nil {
		return nil, result.Error
	}
	var models []model.FollowModel
	for _, follow := range follows {
		models = append(models, &follow)
	}
	return models, nil
}

func (r *repoPrivate) GetFollowing(userID uint) ([]model.FollowModel, error) {
	var follows []Follow
	result := r.db.Where("follower_id = ? AND status = ?", userID, new(FollowAcceptedStatus).ToInt()).Find(&follows)
	if result.Error != nil {
		return nil, result.Error
	}
	var models []model.FollowModel
	for _, follow := range follows {
		models = append(models, &follow)
	}
	return models, nil
}

func (r *repoPrivate) AreAlreadyFollowing(followerID, followedID uint) (bool, error) {
	var follow Follow
	result := r.db.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (r *repoPrivate) IsMyFollow(followerID, followedID uint) (bool, error) {
	var follow Follow
	result := r.db.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Find(&follow)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
