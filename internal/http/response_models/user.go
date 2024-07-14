package response_models

import (
	"go-api/internal/storage/postgres"
	"go-api/pkg/model"
	"time"
)

type GetUserResponseInterface interface {
	GetID() uint
	GetUsername() string
	GetBio() *string
	GetAvatar() *string
	IsPrivateUser() bool
	GetCreatedAt() *time.Time
}

func (u *GetUserResponse) GetID() uint {
	return u.ID
}

func (u *GetUserResponse) GetUsername() string {
	return u.Username
}

func (u *GetUserResponse) GetBio() *string {
	return u.Bio
}

func (u *GetUserResponse) GetAvatar() *string {
	return u.Avatar
}

func (u *GetUserResponse) IsPrivateUser() bool {
	return u.IsPrivate
}

func (u *GetUserResponse) GetCreatedAt() *time.Time {
	return u.CreatedAt
}

type GetUserResponse struct {
	ID        uint
	Username  string
	Bio       *string
	Avatar    *string
	IsPrivate bool
	CreatedAt *time.Time
}

func FormatGetUserResponse(user model.UserModel) GetUserResponseInterface {
	if nil == user {
		return nil
	}
	bio := user.GetBio()
	bioPointer := &bio
	if "" == bio {
		bioPointer = nil
	}

	avatar := user.GetAvatar()
	avatarPointer := &avatar
	if "" == avatar {
		avatarPointer = nil
	}

	createdAt := time.Unix(int64(user.GetCreatedAt()), 0)

	return &GetUserResponse{
		ID:        user.GetID(),
		Username:  user.GetUsername(),
		Bio:       bioPointer,
		Avatar:    avatarPointer,
		IsPrivate: user.IsPrivateUser(),
		CreatedAt: &createdAt,
	}
}

type AdminGetUserResponseInterface interface {
	GetID() uint
	GetUsername() string
	GetRole() string
	GetEmail() string
	GetBio() *string
	GetAvatar() *string
	IsPrivateUser() bool
	GetCreatedAt() *time.Time
}

func (u *AdminGetUserResponse) GetID() uint {
	return u.ID
}

func (u *AdminGetUserResponse) GetUsername() string {
	return u.Username
}

func (u *AdminGetUserResponse) GetRole() string {
	return u.Role
}

func (u *AdminGetUserResponse) GetEmail() string {
	return u.Email
}

func (u *AdminGetUserResponse) GetBio() *string {
	return u.Bio
}

func (u *AdminGetUserResponse) GetAvatar() *string {
	return u.Avatar
}

func (u *AdminGetUserResponse) IsPrivateUser() bool {
	return u.IsPrivate
}

func (u *AdminGetUserResponse) GetCreatedAt() *time.Time {
	return u.CreatedAt
}

type AdminGetUserResponse struct {
	ID        uint
	Username  string
	Role      string
	Email     string
	Bio       *string
	Avatar    *string
	IsPrivate bool
	CreatedAt *time.Time
}

func FormatAdminGetUserResponse(user model.UserModel) AdminGetUserResponseInterface {
	if nil == user {
		return nil
	}
	bio := user.GetBio()
	bioPointer := &bio
	if "" == bio {
		bioPointer = nil
	}

	avatar := user.GetAvatar()
	avatarPointer := &avatar
	if "" == avatar {
		avatarPointer = nil
	}

	createdAt := time.Unix(int64(user.GetCreatedAt()), 0)

	return &AdminGetUserResponse{
		ID:        user.GetID(),
		Username:  user.GetUsername(),
		Role:      user.GetRole(),
		Email:     user.GetEmail(),
		Bio:       bioPointer,
		Avatar:    avatarPointer,
		IsPrivate: user.IsPrivateUser(),
		CreatedAt: &createdAt,
	}
}

type GetOneUserResponse struct {
	ID             uint
	Username       string
	Bio            *string
	Avatar         *string
	IsPrivate      bool
	Email          string
	CreatedAt      *time.Time
	LastDrop       *GetDropResponse
	PinnedDrops    []GetDropResponse
	Groups         []GetGroupResponse
	TotalFollowers int
	TotalFollowed  int
	TotalDrops     int
	CurrentFollow  *GetOneFollowResponse
}

func FormatGetOneUserResponse(
	user model.UserModel,
	lastDrop model.DropModel,
	isLastDropLiked bool,
	pinnedDrops []model.DropModel,
	totalFollowers int,
	totalFollowed int,
	totalDrops int,
	currentFollow model.FollowModel,
	requesterID uint,
) GetOneUserResponse {
	userGroups := user.GetGroups()
	formattedGroups := make([]GetGroupResponse, 0)
	for _, userGroup := range userGroups {
		formattedGroups = append(formattedGroups, FormatGetGroupResponse(userGroup))
	}
	bio := user.GetBio()
	bioPointer := &bio
	if "" == bio {
		bioPointer = nil
	}

	avatar := user.GetAvatar()
	avatarPointer := &avatar
	if "" == avatar {
		avatarPointer = nil
	}

	createdAt := time.Unix(int64(user.GetCreatedAt()), 0)

	formattedPinnedDrops := make([]GetDropResponse, 0)
	var lastDropPointer *GetDropResponse

	if user.IsPrivateUser() &&
		user.GetID() != requesterID &&
		(nil == currentFollow || currentFollow.GetStatus() != new(postgres.FollowAcceptedStatus).ToInt()) {
		lastDrop = nil
		pinnedDrops = nil
	} else {
		for _, drop := range pinnedDrops {
			formattedPinnedDrops = append(formattedPinnedDrops, FormatGetDropResponse(drop, false))
		}

		if nil != lastDrop {
			res := FormatGetDropResponse(lastDrop, isLastDropLiked)
			lastDropPointer = &res
		}
	}

	currentFollowPointer := &GetOneFollowResponse{}
	if nil != currentFollow {
		follow := FormatGetOneFollowResponse(currentFollow)
		currentFollowPointer = &follow
	} else {
		currentFollowPointer = nil
	}

	return GetOneUserResponse{
		ID:             user.GetID(),
		Username:       user.GetUsername(),
		Bio:            bioPointer,
		Avatar:         avatarPointer,
		IsPrivate:      user.IsPrivateUser(),
		Email:          user.GetEmail(),
		CreatedAt:      &createdAt,
		LastDrop:       lastDropPointer,
		PinnedDrops:    formattedPinnedDrops,
		Groups:         formattedGroups,
		TotalFollowers: totalFollowers,
		TotalFollowed:  totalFollowed,
		TotalDrops:     totalDrops,
		CurrentFollow:  currentFollowPointer,
	}
}
