package postgres

import (
	"database/sql"
	"fmt"
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Name         string        `gorm:"not null"`
	Description  string        `gorm:"not null"`
	IsPrivate    bool          `gorm:"default:false"`
	Status       uint          `gorm:"not null;default:1"`
	CreatedByID  uint          `json:"-"`
	CreatedBy    User          `gorm:"foreignKey:CreatedByID;references:ID"`
	GroupMembers []GroupMember `gorm:"foreignKey:GroupID;references:ID"`
	PicturePath  sql.NullString
}

func (g *Group) GetID() uint {
	return g.ID
}

func (g *Group) GetName() string {
	return g.Name
}

func (g *Group) GetDescription() string {
	return g.Description
}

func (g *Group) GetCreatedAt() int {
	return int(g.CreatedAt.Unix())
}

func (g *Group) GetCreatedBy() model.UserModel {
	return &g.CreatedBy
}

func (g *Group) GetPicturePath() sql.NullString {
	return g.PicturePath
}

func (g *Group) GetStatus() uint {
	return g.Status
}

func (g *Group) IsPrivateGroup() bool {
	return g.IsPrivate
}
func (g *Group) GetCreatedByID() uint {
	return g.CreatedByID
}
func (g *Group) GetGroupMembers() []model.GroupMemberModel {
	var members []model.GroupMemberModel
	for _, v := range g.GroupMembers {
		members = append(members, model.GroupMemberModel(&v))
	}

	return members
}

var _ model.GroupModel = (*Group)(nil)

type repoGroupPrivate struct {
	db *gorm.DB
}

var _ model.GroupRepository = (*repoGroupPrivate)(nil)

func (r repoGroupPrivate) Create(name string, description string, isPrivate bool, picturePath string, createdBy model.UserModel) (model.GroupModel, error) {
	group := &Group{
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
		PicturePath: sql.NullString{String: picturePath, Valid: "" != picturePath},
		CreatedBy:   *createdBy.(*User),
	}

	if err := r.db.Create(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}
func NewGroupRepo(db *gorm.DB) model.GroupRepository {
	return &repoGroupPrivate{db: db}
}

func (r repoGroupPrivate) FindAllGroupOwnedByUserId(userId uint) ([]model.GroupModel, error) {
	var groups []*Group

	result := r.db.Preload("GroupMembers").Preload("GroupMembers.Member").Where("created_by_id = ?", userId).Find(&groups)
	if result.Error != nil {
		return nil, result.Error
	}
	models := make([]model.GroupModel, len(groups))
	for i, v := range groups {
		models[i] = model.GroupModel(v)
	}
	return models, nil
}

func (r repoGroupPrivate) Update(args model.FilledGroupPatchParam) (model.GroupModel, error) {
	object := Group{}

	r.db.Where("id = ?", args.ID).First(&object)
	if object.CreatedAt.IsZero() {
		return nil, fmt.Errorf("group with id %d not found", args.ID)
	}

	if args.Name != "" {
		object.Name = args.Name
	}
	if args.Description != "" {
		object.Description = args.Description
	}
	object.IsPrivate = args.IsPrivate
	object.PicturePath = sql.NullString{String: args.Picture, Valid: args.Picture != ""}

	result := r.db.Save(&object)
	return &object, result.Error
}

func (r repoGroupPrivate) GetById(id uint) (model.GroupModel, error) {
	object := Group{}
	object.ID = id

	result := r.db.Preload("CreatedBy").Preload("GroupMembers").Preload("GroupMembers.Member").Find(&object)
	if object.CreatedAt.IsZero() {
		return nil, fmt.Errorf("group with id %d not found", id)
	}

	return &object, result.Error
}

func (r repoGroupPrivate) GetByName(name string) (model.GroupModel, error) {
	object := Group{}

	result := r.db.Preload("GroupMembers").Preload("GroupMembers.Member").Where("name = ?", name).First(&object)
	if object.CreatedAt.IsZero() {
		return nil, fmt.Errorf("group with name %s not found", name)
	}

	return &object, result.Error
}

func (r repoGroupPrivate) Delete(id uint) error {
	//TODO implement me
	panic("implement me")
}

func (r repoGroupPrivate) Search(query string) ([]model.GroupModel, error) {
	var groups []*Group

	result := r.db.Where("LOWER(name) LIKE LOWER(?) AND is_private = false", "%"+query+"%").Find(&groups)
	if result.Error != nil {
		return nil, result.Error
	}

	models := make([]model.GroupModel, len(groups))
	for i, v := range groups {
		models[i] = model.GroupModel(v)
	}

	return models, nil
}

func (r repoGroupPrivate) GetAllGroups() ([]model.GroupModel, error) {
	var foundGroups []*Group
	result := r.db.Preload("GroupMembers").Preload("GroupMembers.Member").Find(&foundGroups)

	models := make([]model.GroupModel, len(foundGroups))
	for i, v := range foundGroups {
		models[i] = model.GroupModel(v)
	}
	return models, result.Error
}
