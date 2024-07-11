package response_models

import (
	"go-api/pkg/custom_type"
	"go-api/pkg/model"
	"time"
)

type GetGroupResponse struct {
	ID           uint
	Name         string
	Description  string
	IsPrivate    bool
	PicturePath  custom_type.NullString
	CreatedAt    *time.Time
	CreatedBy    GetUserResponseInterface            `json:",omitempty"`
	GroupMembers []GetGroupMemberForOneGroupResponse `json:",omitempty"`
}

type GetOneGroupFeedResponse struct {
	ID          uint
	Name        string
	Description string
	IsPrivate   bool
	PicturePath custom_type.NullString
	CreatedAt   *time.Time
	CreatedBy   GetUserResponseInterface `json:",omitempty"`
	GroupDrops  []GetDropResponse
}

type GetGroupMemberForOneGroupResponse struct {
	ID        uint
	Member    GetUserResponseInterface
	Status    uint
	Role      string
	CreatedAt *time.Time
}

type GetSearchGroupResponse struct {
	ID          uint
	Name        string
	Description string
	IsPrivate   bool
	PicturePath custom_type.NullString
	CreatedAt   *time.Time
	CreatedBy   GetUserResponseInterface `json:",omitempty"`
	IsMember    bool
}

func FormatGetOneGroupWithFeed(group model.GroupModel, groupDrops []GetDropResponse) GetOneGroupFeedResponse {
	if nil == group {
		return GetOneGroupFeedResponse{}
	}

	createdAt := time.Unix(int64(group.GetCreatedAt()), 0)

	picturePath := custom_type.NullString{NullString: group.GetPicturePath()}

	return GetOneGroupFeedResponse{
		ID:          group.GetID(),
		Name:        group.GetName(),
		Description: group.GetDescription(),
		IsPrivate:   group.IsPrivateGroup(),
		PicturePath: picturePath,
		CreatedAt:   &createdAt,
		CreatedBy:   FormatGetUserResponse(group.GetCreatedBy()),
		GroupDrops:  groupDrops,
	}
}

func FormatGetGroupResponse(group model.GroupModel) GetGroupResponse {
	if nil == group {
		return GetGroupResponse{}
	}

	createdAt := time.Unix(int64(group.GetCreatedAt()), 0)

	picturePath := custom_type.NullString{NullString: group.GetPicturePath()}

	groupMembers := make([]GetGroupMemberForOneGroupResponse, 0)
	for _, groupMember := range group.GetGroupMembers() {
		groupMembers = append(groupMembers, FormatGetGroupMemberResponseForOneGroup(groupMember))
	}

	return GetGroupResponse{
		ID:           group.GetID(),
		Name:         group.GetName(),
		Description:  group.GetDescription(),
		IsPrivate:    group.IsPrivateGroup(),
		PicturePath:  picturePath,
		CreatedAt:    &createdAt,
		CreatedBy:    FormatGetUserResponse(group.GetCreatedBy()),
		GroupMembers: groupMembers,
	}
}

func FormatGetSearchGroupResponse(group model.GroupModel, isMember bool) GetSearchGroupResponse {
	if nil == group {
		return GetSearchGroupResponse{}
	}

	createdAt := time.Unix(int64(group.GetCreatedAt()), 0)

	picturePath := custom_type.NullString{NullString: group.GetPicturePath()}

	return GetSearchGroupResponse{
		ID:          group.GetID(),
		Name:        group.GetName(),
		Description: group.GetDescription(),
		IsPrivate:   group.IsPrivateGroup(),
		PicturePath: picturePath,
		CreatedAt:   &createdAt,
		CreatedBy:   FormatGetUserResponse(group.GetCreatedBy()),
		IsMember:    isMember,
	}
}

type GetGroupMemberResponse struct {
	ID        uint
	Member    GetUserResponseInterface
	Group     GetGroupResponse
	Status    uint
	Role      string
	CreatedAt *time.Time
}

func FormatGetGroupMemberResponse(groupMember model.GroupMemberModel) GetGroupMemberResponse {
	if nil == groupMember {
		return GetGroupMemberResponse{}
	}

	createdAt := time.Unix(int64(groupMember.GetCreatedAt()), 0)

	return GetGroupMemberResponse{
		ID:        groupMember.GetID(),
		Member:    FormatGetUserResponse(groupMember.GetMember()),
		Group:     FormatGetGroupResponse(groupMember.GetGroup()),
		Status:    groupMember.GetStatus(),
		Role:      groupMember.GetRole(),
		CreatedAt: &createdAt,
	}
}

func FormatGetGroupMemberResponseForOneGroup(groupMember model.GroupMemberModel) GetGroupMemberForOneGroupResponse {
	if nil == groupMember {
		return GetGroupMemberForOneGroupResponse{}
	}

	createdAt := time.Unix(int64(groupMember.GetCreatedAt()), 0)

	return GetGroupMemberForOneGroupResponse{
		ID:        groupMember.GetID(),
		Member:    FormatGetUserResponse(groupMember.GetMember()),
		Status:    groupMember.GetStatus(),
		Role:      groupMember.GetRole(),
		CreatedAt: &createdAt,
	}
}
