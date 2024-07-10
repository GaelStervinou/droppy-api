package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"go-api/pkg/errors2"
	"go-api/pkg/hash"
	"go-api/pkg/model"
	"go-api/pkg/validation"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GoogleID    *string `gorm:"unique"`
	Email       string  `gorm:"unique"`
	Password    string  `gorm:"size:255"`
	Username    string  `gorm:"unique;not null"`
	PhoneNumber string
	Bio         string `gorm:"size:1000"`
	Avatar      string
	VerifyToken string
	Status      int
	IsPrivate   bool `gorm:"default:false"`
	Role        string
	Groups      []Group `gorm:"many2many:group_members;foreignKey:ID;joinForeignKey:MemberID;References:ID;JoinReferences:GroupID"`
}

func (u *User) GetID() uint {
	return u.ID
}

func (u *User) GetGoogleID() *string {
	return u.GoogleID
}

func (u *User) GetEmail() string {
	return u.Email
}
func (u *User) GetPassword() string {
	return u.Password
}
func (u *User) GetUsername() string    { return u.Username }
func (u *User) GetRole() string        { return u.Role }
func (u *User) GetCreatedAt() int      { return int(u.CreatedAt.Unix()) }
func (u *User) GetUpdatedAt() int      { return int(u.UpdatedAt.Unix()) }
func (u *User) GetDeletedAt() int      { return int(u.UpdatedAt.Unix()) }
func (u *User) IsPrivateUser() bool    { return u.IsPrivate }
func (u *User) GetPhoneNumber() string { return u.PhoneNumber }
func (u *User) GetBio() string         { return u.Bio }
func (u *User) GetAvatar() string      { return u.Avatar }
func (u *User) GetGroups() []model.GroupModel {
	var result []model.GroupModel
	for _, userGroup := range u.Groups {
		result = append(result, &userGroup)
	}
	return result
}

var _ model.UserModel = (*User)(nil)

type repoUserPrivate struct {
	db *gorm.DB
}

// Safe checker to know if this file already implements the interface correctly or not
var _ model.UserRepository = (*repoUserPrivate)(nil)

func NewUserRepo(db *gorm.DB) model.UserRepository {
	return &repoUserPrivate{db: db}
}

func (repo *repoUserPrivate) Create(args model.UserCreationParam) (model.UserModel, error) {
	validationError := validation.ValidateUserCreation(args)

	if len(validationError.Fields) > 0 {
		return nil, validationError
	}

	hashedPassword, err := hash.GenerateFromPassword(args.Password)

	if err != nil {
		return nil, err
	}

	userObject := User{
		Email:    args.Email,
		Password: hashedPassword,
		Username: args.Username,
		Role:     args.Role,
		Status:   1,
	}

	result := repo.db.Create(&userObject)

	if result.Error != nil {
		if errors2.IsErrorCode(result.Error, errors2.UniqueViolationErr) {
			return nil, errors.New("email or username already exists")
		}
		return nil, result.Error
	}
	return &userObject, result.Error
}

func (repo *repoUserPrivate) CreateWithGoogle(args model.UserCreationWithGoogleParam) (model.UserModel, error) {
	userObject := User{
		Email:    args.Email,
		GoogleID: &args.GoogleId,
		Status:   1,
		Role:     args.Role,
		Username: args.Username,
	}

	result := repo.db.Create(&userObject)
	return &userObject, result.Error
}

func (repo *repoUserPrivate) Update(args model.UserPatchParam) (model.UserModel, error) {
	validationError := validation.ValidateUserPatch(args)

	if len(validationError.Fields) > 0 {
		return nil, validationError
	}
	userObject := User{}
	repo.db.Where("email = ?", args.Email).First(&userObject)
	if userObject.CreatedAt.IsZero() {
		return nil, errors.New("user not found")
	}

	userObject.Username = args.Username

	result := repo.db.Save(&userObject)
	return &userObject, result.Error
}

func (repo *repoUserPrivate) Delete(id uint) error {
	return repo.db.Delete(&User{}, id).Error
}

func (repo *repoUserPrivate) GetByGoogleAuthId(googleId string) (model.UserModel, error) {
	userObject := User{GoogleID: &googleId}

	result := repo.db.Find(&userObject)
	if userObject.CreatedAt.IsZero() {
		return &userObject, errors.New("user not found")
	}

	return &userObject, result.Error
}

func (repo *repoUserPrivate) GetByEmail(email string) (model.UserModel, error) {
	userObject := User{}
	result := repo.db.Where("email = ?", email).First(&userObject)
	if userObject.CreatedAt.IsZero() {
		return &userObject, errors.New(fmt.Sprintf("user with email %s not found", email))
	}

	return &userObject, result.Error
}

func (repo *repoUserPrivate) GetById(id uint) (model.UserModel, error) {
	userObject := User{}
	userObject.ID = id

	result := repo.db.Preload("Groups").Preload("Groups.CreatedBy").Find(&userObject)
	if userObject.CreatedAt.IsZero() {
		return nil, nil
	}

	return &userObject, result.Error
}

func (repo *repoUserPrivate) GetAll() ([]model.UserModel, error) {
	var foundStudents []*User
	result := repo.db.Find(&foundStudents)

	models := make([]model.UserModel, len(foundStudents))
	for i, v := range foundStudents {
		models[i] = model.UserModel(v)
	}
	return models, result.Error
}
func (repo *repoUserPrivate) CanUserBeFollowed(followedId uint) (bool, error) {
	userObject := User{}
	userObject.ID = followedId

	result := repo.db.Find(&userObject)
	if userObject.CreatedAt.IsZero() {
		return false, errors.New("user not found")
	}

	return userObject.Status == 1, result.Error
}

func (repo *repoUserPrivate) GetUsersFromUserIds(ids []uint) ([]model.UserModel, error) {
	var foundStudents []*User
	result := repo.db.Where("id IN ?", ids).Find(&foundStudents)

	models := make([]model.UserModel, len(foundStudents))
	for i, v := range foundStudents {
		models[i] = model.UserModel(v)
	}
	return models, result.Error
}

func (repo *repoUserPrivate) Search(query string) ([]model.UserModel, error) {
	var foundUsers []*User
	searchParam := "%" + query + "%"
	result := repo.db.Where("LOWER(username) LIKE LOWER(@search)", sql.Named("search", searchParam)).Find(&foundUsers)
	if result.Error != nil {
		return nil, result.Error
	}
	models := make([]model.UserModel, len(foundUsers))
	for i, v := range foundUsers {
		models[i] = model.UserModel(v)
	}
	return models, result.Error
}

func (repo *repoUserPrivate) IsActiveUser(userId uint) (bool, error) {
	userObject := User{}
	userObject.ID = userId

	result := repo.db.Find(&userObject)
	if userObject.CreatedAt.IsZero() {
		return false, errors.New("user not found")
	}

	return userObject.Status == 1, result.Error
}
