package response_models

import (
	"go-api/pkg/model"
	"time"
)

type UserResponse struct {
	ID          uint
	GoogleID    *string
	Email       *string
	Username    string
	Firstname   string
	Lastname    string
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
		Firstname:   user.GetFirstname(),
		Lastname:    user.GetLastname(),
		PhoneNumber: phoneNumberPointer,
		Bio:         bioPointer,
		Avatar:      avatarPointer,
		IsPrivate:   user.IsPrivateUser(),
		Role:        user.GetRole(),
		CreatedAt:   createdAtPointer,
		UpdatedAt:   updatedAtPointer,
	}
}
