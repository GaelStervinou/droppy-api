package model

import (
	"database/sql"
	"mime/multipart"
)

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
	Update(groupID uint, args map[string]interface{}) (GroupModel, error)
	Delete(id uint) error
	Search(query string) ([]GroupModel, error)
	GetAllGroups(page int, pageSize int) ([]GroupModel, error)
	GetAllGroupsCount() (int64, error)
	DeleteGroup(id uint) error
}

type GroupService interface {
	CanCreateGroup(userId uint) (bool, error)
	IsValidGroupCreation(args GroupCreationParam) (bool, error)
	CreateGroup(userId uint, args GroupCreationParam) (GroupModel, error)
	CanUpdateGroup(groupId uint, userId uint) (bool, error)
	IsValidGroupUpdate(groupId uint, args GroupPatchParam) (bool, error)
	PatchGroup(groupId uint, userId uint, args GroupPatchParam) (GroupModel, error)
	GetGroupDrops(groupId uint, requesterID uint) ([]DropModel, error)
	DeleteGroup(groupId uint, userId uint) error
}

type GroupCreationParam struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description"`
	IsPrivate   *bool                 `form:"isPrivate"`
	Picture     *multipart.FileHeader `form:"picture"`
	PicturePath string                `form:"-"`
	Members     []uint                `form:"members"`
}

type GroupPatchParam struct {
	Name        string                `form:"name"  binding:"required"`
	Description string                `form:"description"`
	IsPrivate   *bool                 `form:"isPrivate"`
	Picture     *multipart.FileHeader `form:"picture"`
	PicturePath string                `form:"-"`
}

type FilledGroupPatchParam struct {
	ID          uint
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	IsPrivate   *bool                 `form:"isPrivate" binding:"required"`
	Picture     *multipart.FileHeader `form:"picture"`
	PicturePath string                `form:"-"`
}
