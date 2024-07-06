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
