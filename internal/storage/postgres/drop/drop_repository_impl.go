package drop

import (
	"go-api/internal/storage/postgres/user"
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type Drop struct {
	gorm.Model
	Type               string `gorm:"not null"`
	Content            string `gorm:"not null"`
	Description        string
	CreatedById        uint      `gorm:"not null"`
	CreatedBy          user.User `gorm:"foreignKey:CreatedById;references:ID"`
	Status             uint      `gorm:"not null"`
	DeletedById        uint
	IsPinned           bool `gorm:"default:false"`
	DropNotificationID uint `gorm:"not null"`
	Lat                float64
	Lng                float64
	PicturePath        string
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

type DropStatusActive struct{}

func (d *DropStatusActive) ToInt() uint { return 1 }

type DropStatusDeleted struct{}

func (d *DropStatusDeleted) ToInt() int { return -1 }

type DropStatusBanned struct{}

func (d *DropStatusBanned) ToInt() int { return -2 }

var _ model.DropModel = (*Drop)(nil)

type repoPrivate struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) model.DropRepository {
	return &repoPrivate{db: db}
}

func (r *repoPrivate) Create(
	dropNotificationId uint,
	contentType string,
	content string,
	description string,
	createdById uint,
	status uint,
	isPinned bool,
	picturePath string,
	lat float64,
	lng float64,

) (model.DropModel, error) {
	drop := &Drop{
		Type:               contentType,
		Content:            content,
		Description:        description,
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

func (r *repoPrivate) Delete(dropId uint) error {
	return r.db.Delete(&Drop{}, dropId).Error
}

func (r *repoPrivate) GetDropById(dropId uint) (model.DropModel, error) {
	var drop Drop
	if err := r.db.Preload("CreatedBy").First(&drop, dropId).Error; err != nil {
		return nil, err
	}
	return &drop, nil
}

func (r *repoPrivate) GetUserDrops(userId uint) ([]model.DropModel, error) {
	var drops []Drop
	if err := r.db.Where("created_by_id = ?", userId).Find(&drops).Error; err != nil {
		return nil, err
	}
	var result []model.DropModel
	for _, drop := range drops {
		result = append(result, &drop)
	}
	return result, nil
}

func (r *repoPrivate) GetDropByDropNotificationAndUser(dropNotificationId uint, userId uint) (model.DropModel, error) {
	var drop Drop
	if err := r.db.Where("drop_notification_id = ? AND created_by_id = ?", dropNotificationId, userId).First(&drop).Error; err != nil {
		return nil, err
	}
	return &drop, nil
}

func (r *repoPrivate) GetDropsByUserIdsAndDropNotificationId(userIds []uint, dropNotifId uint) ([]model.DropModel, error) {
	var drops []Drop
	if err := r.db.Preload("CreatedBy").Where("created_by_id IN ? AND drop_notification_id = ?", userIds, dropNotifId).Find(&drops).Error; err != nil {
		return nil, err
	}
	var result []model.DropModel
	for _, drop := range drops {
		result = append(result, &drop)
	}
	return result, nil
}

func (r *repoPrivate) HasUserDropped(dropNotificationId uint, userId uint) (bool, error) {
	var count int64
	if err := r.db.Model(&Drop{}).Where("drop_notification_id = ? AND created_by_id = ?", dropNotificationId, userId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
