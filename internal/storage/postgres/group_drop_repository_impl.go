package postgres

import (
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type GroupDrop struct {
	gorm.Model
	GroupID uint `gorm:"not null"`
	DropID  uint `gorm:"not null"`
	Drop    Drop `gorm:"foreignKey:DropID;references:ID"`
}

func (gd *GroupDrop) GetDropID() uint {
	return gd.DropID
}

func (gd *GroupDrop) GetGroupID() uint {
	return gd.GroupID
}

func (gd *GroupDrop) GetDrop() model.DropModel {
	return &gd.Drop
}

type repoGroupDropPrivate struct {
	db *gorm.DB
}

var _ model.GroupDropRepository = (*repoGroupDropPrivate)(nil)

func NewGroupDropRepo(db *gorm.DB) model.GroupDropRepository {
	return &repoGroupDropPrivate{db: db}
}

func (r repoGroupDropPrivate) Create(dropId uint, groupId uint) (model.GroupDropModel, error) {
	gd := &GroupDrop{
		GroupID: groupId,
		DropID:  dropId,
	}

	if err := r.db.Create(gd).Error; err != nil {
		return nil, err
	}

	return gd, nil
}

func (r repoGroupDropPrivate) Delete(dropId uint, groupId uint) error {
	return r.db.Delete(&GroupDrop{}, dropId, groupId).Error
}

func (r repoGroupDropPrivate) GetByDropId(dropId uint) (model.GroupDropModel, error) {
	var gd GroupDrop
	if err := r.db.First(&gd, dropId).Error; err != nil {
		return nil, err
	}
	return &gd, nil
}

func (r repoGroupDropPrivate) GetByGroupIdAndLastNotificationId(groupId uint, lastNotificationId uint) ([]model.GroupDropModel, error) {
	var gds []GroupDrop
	if err := r.db.
		Joins("JOIN drops ON drops.id = group_drops.drop_id").
		Preload("Drop").
		Preload("Drop.CreatedBy").
		Preload("Drop.Comments").
		Preload("Drop.Comments.CreatedBy").
		Preload("Drop.Comments.Responses").
		Preload("Drop.Comments.Responses.CreatedBy").
		Where("group_id = ? AND drops.drop_notification_id = ?", groupId, lastNotificationId).
		Order("created_at DESC").
		Find(&gds).Error; err != nil {
		return nil, err
	}

	for i := range gds {
		var totalLikes int64
		if err := r.db.Model(&Like{}).Where("drop_id = ?", gds[i].GetDropID()).Count(&totalLikes).Error; err != nil {
			return nil, err
		}
		gds[i].Drop.TotalLikes = int(totalLikes)
	}
	var result []model.GroupDropModel
	for _, gd := range gds {
		result = append(result, &gd)
	}
	return result, nil
}
