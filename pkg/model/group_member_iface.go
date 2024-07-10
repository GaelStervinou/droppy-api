package model

type GroupMemberModel interface {
	GetID() uint
	GetGroupID() uint
	GetGroup() GroupModel
	GetMemberID() uint
	GetMember() UserModel
	GetRole() string
	GetCreatedAt() int
	GetStatus() uint
}

type GroupMemberRepository interface {
	Create(groupID uint, memberID uint, role string, status uint) (GroupMemberModel, error)
	GetByGroupID(groupID uint) ([]GroupMemberModel, error)
	GetByMemberID(memberID uint) ([]GroupMemberModel, error)
	GetByGroupIDAndMemberID(groupID uint, memberID uint) (GroupMemberModel, error)
	UpdateRole(groupID uint, memberID uint, role string) (GroupMemberModel, error)
	UpdateStatus(groupID uint, memberID uint, status uint) (GroupMemberModel, error)
	Delete(groupID uint, memberID uint) error
	IsGroupManager(groupID uint, memberID uint) (bool, error)
	IsGroupMember(groupID uint, memberID uint) (bool, error)
	GetPendingGroupMemberRequests(groupID uint) ([]GroupMemberModel, error)
}

type GroupMemberService interface {
	JoinGroup(currentUserId uint, userId uint, args GroupMemberCreationParam) (GroupMemberModel, error)
	DeleteGroupMember(actionRequesterID uint, groupID uint, memberID uint) error
	AcceptGroupMember(userId uint, groupID uint, memberID uint) (GroupMemberModel, error)
	UpdateGroupMemberRole(requesterId uint, groupID uint, memberID uint, args GroupMemberPatchParam) (GroupMemberModel, error)
	FindAllUserGroups(userId uint) ([]GroupModel, error)
	GetPendingGroupMemberRequests(requesterId uint, groupID uint) ([]GroupMemberModel, error)
	AddUserToGroup(userID uint, groupID uint, requesterID uint) (GroupMemberModel, error)
}

type GroupMemberCreationParam struct {
	GroupID uint   `json:"groupId"`
	Role    string `json:"role"`
}

type GroupMemberPatchParam struct {
	Role string `json:"role"`
}
