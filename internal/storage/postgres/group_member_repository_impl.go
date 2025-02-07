package postgres

import (
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type GroupMember struct {
	gorm.Model
	GroupID  uint   `gorm:"not null"`
	MemberID uint   `gorm:"not null"`
	Status   uint   `gorm:"not null"`
	Role     string `gorm:"not null"`
	Group    Group
	Member   User
}

type GroupMemberStatus interface {
	ToIntGroupMemberStatus() uint
}

type GroupMemberStatusActive struct{}

func (g *GroupMemberStatusActive) ToIntGroupMemberStatus() uint {
	return 1
}

type GroupMemberStatusPending struct{}

func (g *GroupMemberStatusPending) ToIntGroupMemberStatus() uint {
	return 0
}

type GroupMemberRoleManager struct{}

func (g *GroupMemberRoleManager) ToString() string {
	return "manager"
}

type GroupMemberRoleMember struct{}

func (g *GroupMemberRoleMember) ToString() string {
	return "member"
}

func GroupMemberRoles() []string {
	return []string{"manager", "member"}
}

func (g *GroupMember) GetID() uint {
	return g.ID
}

func (g *GroupMember) GetGroupID() uint {
	return g.GroupID
}

func (g *GroupMember) GetGroup() model.GroupModel {
	return &g.Group
}

func (g *GroupMember) GetMemberID() uint {
	return g.MemberID
}

func (g *GroupMember) GetMember() model.UserModel {
	return &g.Member
}

func (g *GroupMember) GetRole() string {
	return g.Role
}

func (g *GroupMember) GetCreatedAt() int {
	return int(g.CreatedAt.Unix())
}

func (g *GroupMember) GetStatus() uint {
	return g.Status
}

var _ model.GroupMemberModel = (*GroupMember)(nil)

type gmRepoPrivate struct {
	db *gorm.DB
}

var _ model.GroupMemberRepository = (*gmRepoPrivate)(nil)

func NewGroupMemberRepo(db *gorm.DB) model.GroupMemberRepository {
	return &gmRepoPrivate{db: db}
}

func (r gmRepoPrivate) Create(groupID uint, memberID uint, role string, status uint) (model.GroupMemberModel, error) {
	groupMember := &GroupMember{
		GroupID:  groupID,
		MemberID: memberID,
		Role:     role,
		Status:   status,
	}

	result := r.db.Create(groupMember)
	if result.Error != nil {
		return nil, result.Error
	}

	return groupMember, nil
}

func (r gmRepoPrivate) GetByGroupID(groupID uint) ([]model.GroupMemberModel, error) {
	var groupMembers []GroupMember
	activeStatus := &GroupMemberStatusActive{}
	result := r.db.Preload("Group").Where("group_id = ? AND status = ?", groupID, activeStatus.ToIntGroupMemberStatus()).Find(&groupMembers)
	if result.Error != nil {
		return nil, result.Error
	}

	var groupMembersModel []model.GroupMemberModel
	for _, groupMember := range groupMembers {
		groupMembersModel = append(groupMembersModel, &groupMember)
	}

	return groupMembersModel, nil
}

func (r gmRepoPrivate) GetByMemberID(memberID uint) ([]model.GroupMemberModel, error) {
	var groupMembers []GroupMember
	activeStatus := &GroupMemberStatusActive{}
	result := r.db.Preload("Group").Where("member_id = ? AND status = ?", memberID, activeStatus.ToIntGroupMemberStatus()).Find(&groupMembers)
	if result.Error != nil {
		return nil, result.Error
	}

	var groupMembersModel []model.GroupMemberModel
	for _, groupMember := range groupMembers {
		groupMembersModel = append(groupMembersModel, &groupMember)
	}

	return groupMembersModel, nil
}

func (r gmRepoPrivate) GetByGroupIDAndMemberID(groupID uint, memberID uint) (model.GroupMemberModel, error) {
	var groupMember GroupMember
	activeStatus := &GroupMemberStatusActive{}
	result := r.db.Preload("Group").Preload("Group.CreatedBy").Preload("Member").Where("group_id = ? AND member_id = ? AND status = ?", groupID, memberID, activeStatus.ToIntGroupMemberStatus()).First(&groupMember)
	if result.Error != nil {
		return nil, result.Error
	}

	return &groupMember, nil
}

func (r gmRepoPrivate) IsGroupMember(groupID uint, memberID uint) (bool, error) {
	var groupMember GroupMember
	result := r.db.Model(&GroupMember{}).Where("group_id = ? AND member_id = ?", groupID, memberID).First(&groupMember)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func (r gmRepoPrivate) UpdateRole(groupID uint, memberID uint, role string) (model.GroupMemberModel, error) {
	result := r.db.Model(&GroupMember{
		MemberID: memberID,
	}).Where("group_id = ? AND member_id = ?", groupID, memberID).Update("role", role)
	if result.Error != nil {
		return nil, result.Error
	}

	finalGroupMember, err := r.GetByGroupIDAndMemberID(groupID, memberID)
	if err != nil {
		return nil, err
	}
	return finalGroupMember, nil
}

func (r gmRepoPrivate) UpdateStatus(groupID uint, memberID uint, status uint) (model.GroupMemberModel, error) {
	var groupMember GroupMember
	result := r.db.Model(&GroupMember{}).Where("group_id = ? AND member_id = ?", groupID, memberID).Update("status", status)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.db.Preload("Group").First(&groupMember)
	if result.Error != nil {
		return nil, result.Error
	}

	return &groupMember, nil
}

func (r gmRepoPrivate) Delete(groupID uint, memberID uint) error {
	result := r.db.Model(&GroupMember{
		MemberID: memberID,
	}).Where("group_id = ? AND member_id = ?", groupID, memberID).Update("status", -1)

	if result.Error != nil {
		return result.Error
	}

	result = r.db.Where("group_id = ? AND member_id = ?", groupID, memberID).Delete(&GroupMember{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r gmRepoPrivate) IsGroupManager(groupID uint, memberID uint) (bool, error) {
	var groupMember GroupMember
	result := r.db.Model(&GroupMember{}).Where("group_id = ? AND member_id = ?", groupID, memberID).First(&groupMember)
	if result.Error != nil {
		return false, result.Error
	}

	role := &GroupMemberRoleManager{}

	return groupMember.GetRole() == role.ToString(), nil
}

func (r gmRepoPrivate) GetPendingGroupMemberRequests(groupID uint) ([]model.GroupMemberModel, error) {
	var groupMembers []GroupMember
	pendingStatus := &GroupMemberStatusPending{}
	result := r.db.Preload("Group").Preload("Member").Preload("Group.CreatedBy").Where("group_id = ? AND status = ?", groupID, pendingStatus.ToIntGroupMemberStatus()).Find(&groupMembers)
	if result.Error != nil {
		return nil, result.Error
	}

	var groupMembersModel []model.GroupMemberModel
	for _, groupMember := range groupMembers {
		groupMembersModel = append(groupMembersModel, &groupMember)
	}

	return groupMembersModel, nil
}

func (r gmRepoPrivate) DeleteGroupMembers(groupID uint) error {
	result := r.db.Where("group_id = ?", groupID).Delete(&GroupMember{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r gmRepoPrivate) IsUserInGroups(groupIds []uint, memberID uint) (bool, error) {
	var groupMembers []GroupMember
	activeStatus := &GroupMemberStatusActive{}
	result := r.db.Preload("Group").Preload("Group.CreatedBy").Preload("Member").Where("group_id IN ? AND member_id = ? AND status = ?", groupIds, memberID, activeStatus.ToIntGroupMemberStatus()).Find(&groupMembers)
	if result.Error != nil {
		return false, result.Error
	}

	return len(groupMembers) > 0, nil
}
