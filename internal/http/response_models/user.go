package response_models

import (
	"go-api/pkg/model"
	"time"
)

type GetUserResponse struct {
	ID        uint
	Username  string
	Bio       *string
	Avatar    *string
	IsPrivate bool
	CreatedAt *time.Time
}

type UserResponse struct {
	ID          uint
	GoogleID    *string
	Email       *string
	Username    string
	PhoneNumber *string
	Bio         *string
	Avatar      *string
	IsPrivate   bool
	Role        string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (u *UserResponse) HidePersonalInfo() {
	u.Email = nil
	u.PhoneNumber = nil
	u.GoogleID = nil
	u.CreatedAt = nil
	u.UpdatedAt = nil
}

func FormatGetUserResponse(user model.UserModel) GetUserResponse {
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

	return GetUserResponse{
		ID:        user.GetID(),
		Username:  user.GetUsername(),
		Bio:       bioPointer,
		Avatar:    avatarPointer,
		IsPrivate: user.IsPrivateUser(),
		CreatedAt: &createdAt,
	}
}

func FormatUserFromModel(user model.UserModel) UserResponse {
	email := user.GetEmail()
	emailPointer := &email
	phoneNumber := user.GetPhoneNumber()
	phoneNumberPointer := &phoneNumber
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
	updatedAt := time.Unix(int64(user.GetUpdatedAt()), 0)
	createdAtPointer := &createdAt
	updatedAtPointer := &updatedAt

	return UserResponse{
		ID:          user.GetID(),
		GoogleID:    user.GetGoogleID(),
		Email:       emailPointer,
		Username:    user.GetUsername(),
		PhoneNumber: phoneNumberPointer,
		Bio:         bioPointer,
		Avatar:      avatarPointer,
		IsPrivate:   user.IsPrivateUser(),
		Role:        user.GetRole(),
		CreatedAt:   createdAtPointer,
		UpdatedAt:   updatedAtPointer,
	}
}
