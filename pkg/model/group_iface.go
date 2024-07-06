package model

import "database/sql"

type GroupModel interface {
	GetID() uint
	GetName() string
	GetDescription() string
	GetCreatedAt() int
	GetCreatedBy() UserModel
	IsPrivateGroup() bool
	GetPicturePath() sql.NullString
}

type GroupRepository interface {
	Create(name string, description string, isPrivate bool, picturePath string, createdBy UserModel) (GroupModel, error)
	FindAllByUserId(userId uint) ([]GroupModel, error)
	GetById(userId uint, id uint) (GroupModel, error)
	GetByName(name string) (GroupModel, error)
	Delete(id uint) error
}

type GroupService interface {
	CanCreateGroup(userId uint, args GroupCreationParam) (bool, error)
	IsValidGroupCreation(args GroupCreationParam) (bool, error)
	CreateGroup(userId uint, args GroupCreationParam) (GroupModel, error)
}

type GroupCreationParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"isPrivate"`
	Picture     string `json:"picture"`
}
