package postgres

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

type repoDropNotifPrivate struct {
	db *gorm.DB
}

func NewDropNotifRepo(db *gorm.DB) model.DropNotificationRepository {
	return &repoDropNotifPrivate{db: db}
}

func (r *repoDropNotifPrivate) Create(notificationType string) (model.DropNotificationModel, error) {
	notification := &DropNotification{
		Type: notificationType,
	}
	r.db.Create(notification)
	return notification, nil
}

func (r *repoDropNotifPrivate) GetNotificationByID(notificationId uint) (model.DropNotificationModel, error) {
	var notification DropNotification
	r.db.First(&notification, notificationId)
	return &notification, nil
}

func (r *repoDropNotifPrivate) GetCurrentDropNotification() (model.DropNotificationModel, error) {
	var notification DropNotification
	r.db.Order("id desc").First(&notification)
	return &notification, nil
}
