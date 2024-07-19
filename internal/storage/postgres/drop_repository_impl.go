package postgres

import (
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type Drop struct {
	gorm.Model
	Type               string `gorm:"not null"`
	ContentTitle       string `gorm:"not null"`
	ContentSubtitle    string `gorm:"default:null"`
	Location           string
	Content            string `gorm:"not null"`
	ContentPicturePath string `gorm:"not null"`
	Description        string
	CreatedById        uint `gorm:"not null;index:idx_drop_notification_created_by"`
	CreatedBy          User `gorm:"foreignKey:CreatedById;references:ID"`
	Status             uint `gorm:"not null"`
	DeletedById        uint
	IsPinned           bool `gorm:"default:false"`
	DropNotificationID uint `gorm:"not null;index:idx_drop_notification_created_by"`
	Lat                float64
	Lng                float64
	PicturePath        string
	Comments           []Comment `gorm:"foreignKey:DropId;references:ID"`
	TotalLikes         int       `gorm:"-"`
}

func (d *Drop) GetID() uint { return d.ID }

func (d *Drop) GetDropNotificationID() uint { return d.DropNotificationID }

func (d *Drop) GetType() string { return d.Type }

func (d *Drop) GetContent() string { return d.Content }

func (d *Drop) GetDescription() string { return d.Description }

func (d *Drop) GetCreatedById() uint { return d.CreatedById }

func (d *Drop) GetStatus() uint { return d.Status }

func (d *Drop) GetDeletedById() uint { return d.DeletedById }

func (d *Drop) GetIsPinned() bool { return d.IsPinned }

func (d *Drop) GetLat() float64 { return d.Lat }

func (d *Drop) GetLng() float64 { return d.Lng }

func (d *Drop) GetPicturePath() string { return d.PicturePath }

func (d *Drop) GetCreatedAt() int { return int(d.CreatedAt.Unix()) }

func (d *Drop) GetCreatedBy() model.UserModel { return &d.CreatedBy }

func (d *Drop) GetComments() []model.CommentModel {
	var result []model.CommentModel
	for _, comment := range d.Comments {
		result = append(result, &comment)
	}
	return result
}

func (d *Drop) GetTotalLikes() int { return d.TotalLikes }

func (d *Drop) GetContentTitle() string { return d.ContentTitle }

func (d *Drop) GetContentSubtitle() string { return d.ContentSubtitle }

func (d *Drop) GetContentPicturePath() string { return d.ContentPicturePath }

func (d *Drop) GetLocation() string { return d.Location }

type DropStatusActive struct{}

func (d *DropStatusActive) ToInt() uint { return 1 }

type DropStatusDeleted struct{}

func (d *DropStatusDeleted) ToInt() int { return -1 }

type DropStatusBanned struct{}

func (d *DropStatusBanned) ToInt() int { return -2 }

var _ model.DropModel = (*Drop)(nil)

type repoDropPrivate struct {
	db *gorm.DB
}

func NewDropRepo(db *gorm.DB) model.DropRepository {
	return &repoDropPrivate{db: db}
}

func (r *repoDropPrivate) Create(
	dropNotificationId uint,
	contentType string,
	content string,
	description string,
	contentPicturePath string,
	contentTitle string,
	contentSubtitle string,
	createdById uint,
	status uint,
	isPinned bool,
	picturePath string,
	lat float64,
	lng float64,
	location string,

) (model.DropModel, error) {
	drop := &Drop{
		Type:               contentType,
		Content:            content,
		Description:        description,
		ContentPicturePath: contentPicturePath,
		ContentTitle:       contentTitle,
		ContentSubtitle:    contentSubtitle,
		Location:           location,
		CreatedById:        createdById,
		Status:             status,
		IsPinned:           isPinned,
		DropNotificationID: dropNotificationId,
		PicturePath:        picturePath,
		Lat:                lat,
		Lng:                lng,
	}
	if err := r.db.Create(drop).Error; err != nil {
		return nil, err
	}
	return drop, nil
}

func (r *repoDropPrivate) Delete(dropId uint) error {
	return r.db.Delete(&Drop{}, dropId).Error
}

func (r *repoDropPrivate) CountUserDrops(userId uint) int {
	var count int64
	r.db.Model(&Drop{}).Where("created_by_id = ?", userId).Count(&count)
	return int(count)
}

func (r *repoDropPrivate) GetDropGroups(dropId uint) ([]model.GroupModel, error) {
	var groups []Group
	if err := r.db.Where("id = ?", dropId).Preload("Group").Find(&groups).Error; err != nil {
		return nil, err
	}
	var result []model.GroupModel
	for _, group := range groups {
		result = append(result, &group)
	}
	return result, nil
}

func (r *repoDropPrivate) CountGroupDrops(groupId uint) int {
	var count int64
	r.db.Model(&Drop{}).Where("drop_notification_id = ?", groupId).Count(&count)
	return int(count)
}

func (r *repoDropPrivate) GetDropById(dropId uint) (model.DropModel, error) {
	var drop Drop
	if err := r.db.
		Preload("CreatedBy").
		Preload("Comments").
		Preload("Comments.CreatedBy").
		Preload("Comments.Responses").
		Preload("Comments.Responses.CreatedBy").
		First(&drop, dropId).Error; err != nil {
		return nil, err
	}
	var totalLikes int64
	if err := r.db.Model(&Like{}).Where("drop_id = ?", drop.ID).Count(&totalLikes).Error; err != nil {
		return nil, err
	}
	drop.TotalLikes = int(totalLikes)
	return &drop, nil
}

func (r *repoDropPrivate) GetUserDrops(userId uint) ([]model.DropModel, error) {
	var drops []Drop
	if err := r.db.Preload("CreatedBy").Preload("Comments").Where("created_by_id = ?", userId).Find(&drops).Error; err != nil {
		return nil, err
	}
	var result []model.DropModel
	for _, drop := range drops {
		result = append(result, &drop)
	}
	return result, nil
}

func (r *repoDropPrivate) GetDropByDropNotificationAndUser(dropNotificationId uint, userId uint) (model.DropModel, error) {
	var drop Drop
	if err := r.db.Preload("CreatedBy").Preload("Comments").Where("drop_notification_id = ? AND created_by_id = ?", dropNotificationId, userId).First(&drop).Error; err != nil {
		return nil, err
	}
	return &drop, nil
}

func (r *repoDropPrivate) DropExists(dropId uint) (bool, error) {
	var count int64
	if err := r.db.Model(&Drop{}).Where("id = ?", dropId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repoDropPrivate) GetDropsByUserIdsAndDropNotificationId(userIds []uint, dropNotifId uint) ([]model.DropModel, error) {
	var drops []Drop
	if err := r.db.
		Preload("CreatedBy").
		Preload("Comments").
		Preload("Comments.CreatedBy").
		Preload("Comments.Responses").
		Preload("Comments.Responses.CreatedBy").
		Where("created_by_id IN ? AND drop_notification_id = ?", userIds, dropNotifId).
		Order("created_at desc").
		Find(&drops).Error; err != nil {
		return nil, err
	}

	for i := range drops {
		var totalLikes int64
		if err := r.db.Model(&Like{}).Where("drop_id = ?", drops[i].ID).Count(&totalLikes).Error; err != nil {
			return nil, err
		}
		drops[i].TotalLikes = int(totalLikes)
	}

	var result []model.DropModel
	for _, drop := range drops {
		result = append(result, &drop)
	}
	return result, nil
}

func (r *repoDropPrivate) HasUserDropped(dropNotificationId uint, userId uint) (bool, error) {
	var count int64
	if err := r.db.Model(&Drop{}).Where("drop_notification_id = ? AND created_by_id = ?", dropNotificationId, userId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repoDropPrivate) GetUserPinnedDrops(userId uint) ([]model.DropModel, error) {
	var drops []Drop
	if err := r.db.
		Preload("CreatedBy").
		Preload("Comments").
		Preload("Comments.CreatedBy").
		Preload("Comments.Responses").
		Preload("Comments.Responses.CreatedBy").
		Where("created_by_id = ? AND is_pinned = ?", userId, true).
		Order("created_at desc").
		Find(&drops).
		Error; err != nil {
		return nil, err
	}
	var result []model.DropModel
	for _, drop := range drops {
		result = append(result, &drop)
	}
	return result, nil
}

func (r *repoDropPrivate) GetUserLastDrop(userId uint, lastNotifID uint) (model.DropModel, error) {
	var drop Drop
	if err := r.db.
		Preload("CreatedBy").
		Preload("Comments").
		Preload("Comments.Responses").
		Preload("Comments.Responses.CreatedBy").
		Where("created_by_id = ? AND drop_notification_id = ?", userId, lastNotifID).
		Order("created_at desc").
		First(&drop).Error; err != nil {
		return nil, err
	}

	return &drop, nil
}

func (r *repoDropPrivate) GetAllDrops(page int, pageSize int) ([]model.DropModel, error) {
	var drops []Drop
	offset := (page - 1) * pageSize
	if err := r.db.
		Preload("CreatedBy").
		Preload("Comments").
		Preload("Comments.CreatedBy").
		Preload("Comments.Responses").
		Preload("Comments.Responses.CreatedBy").
		Order("id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&drops).Error; err != nil {
		return nil, err
	}

	for i := range drops {
		var totalLikes int64
		if err := r.db.Model(&Like{}).Where("drop_id = ?", drops[i].ID).Count(&totalLikes).Error; err != nil {
			return nil, err
		}
		drops[i].TotalLikes = int(totalLikes)
	}

	var result []model.DropModel
	for _, drop := range drops {
		result = append(result, &drop)
	}
	return result, nil
}

func (r *repoDropPrivate) GetAllDropsCount() (int64, error) {
	var count int64
	if err := r.db.Model(&Drop{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repoDropPrivate) Update(dropId uint, updates map[string]interface{}) (model.DropModel, error) {
	var drop Drop
	if err := r.db.Model(&Drop{}).Where("id = ?", dropId).Updates(updates).First(&drop).Error; err != nil {
		return nil, err
	}
	return &drop, nil
}
