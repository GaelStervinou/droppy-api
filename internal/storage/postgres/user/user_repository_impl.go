package user

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"go-api/pkg/model"
	"gorm.io/gorm"
	"strings"
)

type User struct {
	gorm.Model
	GoogleID    *string `gorm:"unique"`
	Email       string  `gorm:"unique"`
	Password    string  `json:"-"`
	Username    string  `gorm:"unique;not null"`
	Firstname   string
	Lastname    string
	PhoneNumber string
	Bio         string `gorm:"size:1000"`
	Avatar      string
	VerifyToken string
	Status      int
	IsPrivate   bool        `gorm:"default:false"`
	Roles       StringSlice `gorm:"type:VARCHAR(255)"`
}
type StringSlice []string

func (s *StringSlice) Scan(src any) error {
	if strings.Contains(src.(string), ",") {
		*s = strings.Split(src.(string), ",")
		return nil
	} else {
		*s = []string{src.(string)}
		return nil
	}
}
func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return strings.Join(s, ","), nil
}

func (u *User) GetID() uint {
	return u.ID
}

func (u *User) GetGoogleID() string {
	return *u.GoogleID
}

func (u *User) GetEmail() string {
	return u.Email
}
func (u *User) GetPassword() string {
	return u.Password
}
func (u *User) GetFirstname() string {
	return u.Firstname
}
func (u *User) GetLastname() string {
	return u.Lastname
}
func (u *User) GetUsername() string { return u.Username }
func (u *User) GetRoles() []string  { return u.Roles }
func (u *User) GetCreatedAt() int   { return int(u.CreatedAt.Unix()) }
func (u *User) GetUpdatedAt() int   { return int(u.UpdatedAt.Unix()) }
func (u *User) GetDeletedAt() int   { return int(u.UpdatedAt.Unix()) }

var _ model.UserModel = (*User)(nil)

type repoPrivate struct {
	db *gorm.DB
}

// Safe checker to know if this file already implements the interface correctly or not
var _ model.UserRepository = (*repoPrivate)(nil)

func NewRepo(db *gorm.DB) model.UserRepository {
	return &repoPrivate{db: db}
}

func (repo *repoPrivate) Create(args model.UserCreationParam) (model.UserModel, error) {
	userObject := User{
		Firstname: args.Firstname,
		Lastname:  args.Lastname,
		Email:     args.Email,
		Password:  args.Password,
		Roles:     args.Roles,
	}

	result := repo.db.Create(&userObject)
	return &userObject, result.Error
}

func (repo *repoPrivate) CreateWithGoogle(args model.UserCreationWithGoogleParam) (model.UserModel, error) {
	userObject := User{
		Firstname: args.Firstname,
		Lastname:  args.Lastname,
		Email:     args.Email,
		GoogleID:  &args.GoogleId,
		Roles:     args.Roles,
	}

	result := repo.db.Create(&userObject)
	return &userObject, result.Error
}

func (repo *repoPrivate) Update(user model.UserModel) (model.UserModel, error) {
	return user, repo.db.Save(user).Error
}

func (repo *repoPrivate) Delete(id uint) error {
	return repo.db.Delete(&User{}, id).Error
}

func (repo *repoPrivate) GetByGoogleAuthId(googleId string) (model.UserModel, error) {
	userObject := User{GoogleID: &googleId}

	result := repo.db.Find(&userObject)
	if userObject.CreatedAt.IsZero() {
		return &userObject, errors.New("user not found")
	}

	return &userObject, result.Error
}

func (repo *repoPrivate) GetByEmail(email string) (model.UserModel, error) {
	userObject := User{Email: email}

	result := repo.db.Find(&userObject)
	if userObject.CreatedAt.IsZero() {
		return &userObject, errors.New(fmt.Sprintf("user with email %s not found", email))
	}

	return &userObject, result.Error
}

func (repo *repoPrivate) GetById(id uint) (model.UserModel, error) {
	userObject := User{}
	userObject.ID = id

	result := repo.db.Find(&userObject)
	if userObject.CreatedAt.IsZero() {
		return &userObject, errors.New("user not found")
	}

	return &userObject, result.Error
}

func (repo *repoPrivate) GetAll() ([]model.UserModel, error) {
	var foundStudents []*User
	result := repo.db.Find(&foundStudents)

	models := make([]model.UserModel, len(foundStudents))
	for i, v := range foundStudents {
		models[i] = model.UserModel(v)
	}
	return models, result.Error
}
