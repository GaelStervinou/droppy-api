package model

type UserModel interface {
	GetID() uint
	GetPassword() string
	GetRole() string
	GetEmail() string
	GetGoogleID() *string
	GetUsername() string
	GetPhoneNumber() string
	GetBio() string
	GetAvatar() string
	IsPrivateUser() bool
	GetCreatedAt() int
	GetUpdatedAt() int
	GetGroups() []GroupModel
}

type UserRepository interface {
	GetByGoogleAuthId(googleID string) (UserModel, error)
	GetByEmail(email string) (UserModel, error)
	GetById(id uint) (UserModel, error)
	Create(args UserCreationParam) (UserModel, error)
	CreateWithGoogle(args UserCreationWithGoogleParam) (UserModel, error)
	Update(args UserPatchParam) (UserModel, error)
	Delete(id uint) error
	GetAll() ([]UserModel, error)
	CanUserBeFollowed(followedId uint) (bool, error)
	GetUsersFromUserIds(userIds []uint) ([]UserModel, error)
	Search(query string) ([]UserModel, error)
	IsActiveUser(userId uint) (bool, error)
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
	Email    string
	Username string
}

type LoginParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
