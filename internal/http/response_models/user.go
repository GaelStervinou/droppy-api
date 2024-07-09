package response_models

import (
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
	GetPhoneNumber() string
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

func (u *AdminGetUserResponse) GetPhoneNumber() string {
	return u.PhoneNumber
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
	ID          uint
	Username    string
	Role        string
	PhoneNumber string
	Email       string
	Bio         *string
	Avatar      *string
	IsPrivate   bool
	CreatedAt   *time.Time
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
		ID:          user.GetID(),
		Username:    user.GetUsername(),
		Role:        user.GetRole(),
		PhoneNumber: user.GetPhoneNumber(),
		Email:       user.GetEmail(),
		Bio:         bioPointer,
		Avatar:      avatarPointer,
		IsPrivate:   user.IsPrivateUser(),
		CreatedAt:   &createdAt,
	}
}

type GetOneUserResponse struct {
	ID             uint
	Username       string
	Bio            *string
	Avatar         *string
	IsPrivate      bool
	CreatedAt      *time.Time
	LastDrop       GetDropResponse
	PinnedDrops    []GetDropResponse
	TotalFollowers int
	TotalFollowed  int
}

func FormatGetOneUserResponse(
	user model.UserModel,
	lastDrop model.DropModel,
	isLastDropLiked bool,
	pinnedDrops []model.DropModel,
	totalFollowers int,
	totalFollowed int,
) GetOneUserResponse {
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
	for _, drop := range pinnedDrops {
		formattedPinnedDrops = append(formattedPinnedDrops, FormatGetDropResponse(drop, false))
	}

	return GetOneUserResponse{
		ID:             user.GetID(),
		Username:       user.GetUsername(),
		Bio:            bioPointer,
		Avatar:         avatarPointer,
		IsPrivate:      user.IsPrivateUser(),
		CreatedAt:      &createdAt,
		LastDrop:       FormatGetDropResponse(lastDrop, isLastDropLiked),
		PinnedDrops:    formattedPinnedDrops,
		TotalFollowers: totalFollowers,
		TotalFollowed:  totalFollowed,
	}
}
