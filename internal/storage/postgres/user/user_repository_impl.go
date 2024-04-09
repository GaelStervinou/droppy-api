package user

import (
	"context"
	"errors"
	"fmt"
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GoogleID  *string `gorm:"unique"`
	Email     string  `gorm:"unique"`
	Password  string  `json:"-"`
	Firstname string
	Lastname  string
	Role      string
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

func (u *User) GetRole() string {
	return u.Role
}

func (u *User) GetCreatedAt() int {
	return int(u.CreatedAt.Unix())
}

func (u *User) GetUpdatedAt() int {
	return int(u.UpdatedAt.Unix())
}

var _ model.UserModel = (*User)(nil)

type repoPrivate struct {
	db *gorm.DB
}

// Safe checker to know if this file already implements the interface correctly or not
var _ model.UserRepository = (*repoPrivate)(nil)

func NewRepo(db *gorm.DB) model.UserRepository {
	return &repoPrivate{db: db}
}

func (repo *repoPrivate) Create(ctx context.Context, args model.UserCreationParam) (model.UserModel, error) {
	userObject := User{
		Firstname: args.Firstname,
		Lastname:  args.Lastname,
		Email:     args.Email,
		Password:  args.Password,
		Role:      args.Role,
	}

	result := repo.db.Create(&userObject)
	return &userObject, result.Error
}

func (repo *repoPrivate) CreateWithGoogle(ctx context.Context, args model.UserCreationWithGoogleParam) (model.UserModel, error) {
	userObject := User{
		Firstname: args.Firstname,
		Lastname:  args.Lastname,
		Email:     args.Email,
		GoogleID:  &args.GoogleId,
		Role:      args.Role,
	}

	result := repo.db.Create(&userObject)
	return &userObject, result.Error
}

func (repo *repoPrivate) Update(ctx context.Context, user model.UserModel) (model.UserModel, error) {
	return user, repo.db.Save(user).Error
}

func (repo *repoPrivate) Delete(ctx context.Context, id uint) error {
	return repo.db.Delete(&User{}, id).Error
}

func (repo *repoPrivate) GetByGoogleAuthId(ctx context.Context, googleId string) (model.UserModel, error) {
	userObject := User{GoogleID: &googleId}

	result := repo.db.Find(&userObject)
	if userObject.CreatedAt.IsZero() {
		return &userObject, errors.New("user not found")
	}

	return &userObject, result.Error
}

func (repo *repoPrivate) GetByEmail(ctx context.Context, email string) (model.UserModel, error) {
	userObject := User{Email: email}

	result := repo.db.Find(&userObject)
	if userObject.CreatedAt.IsZero() {
		return &userObject, errors.New(fmt.Sprintf("user with email %s not found", email))
	}

	return &userObject, result.Error
}

func (repo *repoPrivate) GetById(ctx context.Context, id uint) (model.UserModel, error) {
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
