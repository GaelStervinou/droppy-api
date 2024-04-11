package model

type UserModel interface {
	GetID() uint
	GetPassword() string
	GetRoles() []string
}

type UserRepository interface {
	GetByGoogleAuthId(googleID string) (UserModel, error)
	GetByEmail(email string) (UserModel, error)
	GetById(id uint) (UserModel, error)
	Create(args UserCreationParam) (UserModel, error)
	CreateWithGoogle(args UserCreationWithGoogleParam) (UserModel, error)
	Update(user UserModel) (UserModel, error)
	Delete(id uint) error
	GetAll() ([]UserModel, error)
}

type UserCreationParam struct {
	Firstname string
	Lastname  string
	Email     string
	Password  string
	Username  string
	Roles     []string
}

type UserCreationWithGoogleParam struct {
	Firstname string
	Lastname  string
	Email     string
	Username  string
	GoogleId  string
	Roles     []string
}
