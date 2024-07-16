package model

type DropNotificationModel interface {
	GetID() uint
	GetType() string
	GetCreatedAt() string
}

type DropNotificationRepository interface {
	Create(notificationType string) (DropNotificationModel, error)
	GetNotificationByID(notificationId uint) (DropNotificationModel, error)
	GetCurrentDropNotification() (DropNotificationModel, error)
}

type ScheduleDropParam struct {
	Type string `json:"type" binding:"required"`
}
