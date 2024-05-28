package model

type DropNotificationModel interface {
	GetID() uint
	GetType() string
	GetCreatedAt() string
}

type DropNotificationRepository interface {
	Create(dropId, createdById uint, notificationType string) (DropNotificationModel, error)
	GetNotificationByID(notificationId uint) (DropNotificationModel, error)
	GetCurrentDropNotification() (DropNotificationModel, error)
}
