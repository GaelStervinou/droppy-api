package drop_notification

import (
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type DropNotification struct {
	gorm.Model
	Type string
}

func (d *DropNotification) GetID() uint { return d.ID }

func (d *DropNotification) GetType() string { return d.Type }

func (d *DropNotification) GetCreatedAt() string { return d.CreatedAt.String() }

type repoPrivate struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) model.DropNotificationRepository {
	return &repoPrivate{db: db}
}

func (r *repoPrivate) Create(dropId, createdById uint, notificationType string) (model.DropNotificationModel, error) {
	notification := &DropNotification{
		Type: notificationType,
	}
	r.db.Create(notification)
	return notification, nil
}

func (r *repoPrivate) GetNotificationByID(notificationId uint) (model.DropNotificationModel, error) {
	var notification DropNotification
	r.db.First(&notification, notificationId)
	return &notification, nil
}

func (r *repoPrivate) GetCurrentDropNotification() (model.DropNotificationModel, error) {
	var notification DropNotification
	r.db.Last(&notification)
	return &notification, nil
}
