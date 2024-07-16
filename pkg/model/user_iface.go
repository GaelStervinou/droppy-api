package model

import "mime/multipart"

type UserModel interface {
	GetID() uint
	GetPassword() string
	GetRole() string
	GetEmail() string
	GetFirebaseUID() string
	GetUsername() string
	GetBio() string
	GetStatus() int
	GetAvatar() string
	IsPrivateUser() bool
	GetCreatedAt() int
	GetUpdatedAt() int
	GetGroups() []GroupModel
	GetFCMToken() string
}

type UserRepository interface {
	GetByFirebaseUid(googleID string) (UserModel, error)
	GetByEmail(email string) (UserModel, error)
	GetById(id uint) (UserModel, error)
	Create(args UserCreationParam) (UserModel, error)
	CreateWithGoogle(args UserCreationWithGoogleParam) (UserModel, error)
	Update(userID uint, args map[string]interface{}) (UserModel, error)
	Delete(id uint) error
	GetAll() ([]UserModel, error)
	CanUserBeFollowed(followedId uint) (bool, error)
	GetUsersFromUserIds(userIds []uint) ([]UserModel, error)
	Search(query string) ([]UserModel, error)
	IsActiveUser(userId uint) (bool, error)
	GetAllFCMTokens() ([]string, error)
	BanUser(userId uint) (UserModel, error)
	UnbanUser(userId uint) (UserModel, error)
	UpdateByAdmin(userId uint, args AdminUpdateUserRequest) (UserModel, error)
}

type UserService interface {
	UpdateUser(userId uint, args UserPatchParam) (UserModel, error)
}

type UserCreationParam struct {
	Email    string
	Password string
	Username string
	Role     string
}

type UserCreationWithGoogleParam struct {
	Email    string
	Username string
	GoogleId string
	Role     string
}

type UserPatchParam struct {
	Bio         string                `form:"bio"`
	Username    string                `form:"username"`
	IsPrivate   *bool                 `form:"isPrivate"`
	Picture     *multipart.FileHeader `form:"picture"`
	PicturePath string                `form:"-"`
}

type LoginParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FcmToken string `json:"fcmToken"`
}

type AdminUpdateUserRequest struct {
	Role     string `json:"role" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Status   int    `json:"status"`
}
