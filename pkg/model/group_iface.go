package model

import "database/sql"

type GroupModel interface {
	GetID() uint
	GetName() string
	GetDescription() string
	GetCreatedAt() int
	GetCreatedByID() uint
	GetCreatedBy() UserModel
	IsPrivateGroup() bool
	GetPicturePath() sql.NullString
	GetGroupMembers() []GroupMemberModel
}

type GroupRepository interface {
	Create(name string, description string, isPrivate bool, picturePath string, createdBy UserModel) (GroupModel, error)
	FindAllGroupOwnedByUserId(userId uint) ([]GroupModel, error)
	GetById(id uint) (GroupModel, error)
	GetByName(name string) (GroupModel, error)
	Update(args FilledGroupPatchParam) (GroupModel, error)
	Delete(id uint) error
	Search(query string) ([]GroupModel, error)
	GetAllGroups() ([]GroupModel, error)
}

type GroupService interface {
	CanCreateGroup(userId uint) (bool, error)
	IsValidGroupCreation(args GroupCreationParam) (bool, error)
	CreateGroup(userId uint, args GroupCreationParam) (GroupModel, error)
	CanUpdateGroup(groupId uint, userId uint) (bool, error)
	IsValidGroupUpdate(groupId uint, args GroupPatchParam) (bool, error)
	PatchGroup(groupId uint, userId uint, args GroupPatchParam) (GroupModel, error)
}

type GroupCreationParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"isPrivate"`
	Picture     string `json:"picture"`
}

type GroupPatchParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"isPrivate"`
	Picture     string `json:"picture"`
}

type FilledGroupPatchParam struct {
	ID          uint
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"isPrivate"`
	Picture     string `json:"picture"`
}
