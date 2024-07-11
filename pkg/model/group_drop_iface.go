package model

type GroupDropModel interface {
	GetDropID() uint
	GetGroupID() uint
	GetDrop() DropModel
}

type GroupDropRepository interface {
	Create(dropId uint, groupId uint) (GroupDropModel, error)
	Delete(dropId uint, groupId uint) error
	GetByDropId(dropId uint) (GroupDropModel, error)
	GetByGroupIdAndLastNotificationId(groupId uint, lastNotificationId uint) ([]GroupDropModel, error)
}
