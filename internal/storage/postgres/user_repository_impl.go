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
	FirebaseUID string
	Email       string `gorm:"unique"`
	Password    string `gorm:"size:255"`
	Username    string `gorm:"unique;not null"`
	Bio         string `gorm:"size:1000"`
	Avatar      string
	VerifyToken string
	Status      int
	IsPrivate   bool `gorm:"default:false"`
	Role        string
	Groups      []Group `gorm:"many2many:group_members;foreignKey:ID;joinForeignKey:MemberID;References:ID;JoinReferences:GroupID"`
	FCMToken    string
}

func (u *User) GetID() uint {
	return u.ID
}

func (u *User) GetFirebaseUID() string {
	return u.FirebaseUID
}

func (u *User) GetEmail() string {
	return u.Email
}
func (u *User) GetPassword() string {
	return u.Password
}
func (u *User) GetUsername() string { return u.Username }
func (u *User) GetRole() string     { return u.Role }
func (u *User) GetCreatedAt() int   { return int(u.CreatedAt.Unix()) }
func (u *User) GetUpdatedAt() int   { return int(u.UpdatedAt.Unix()) }
func (u *User) GetDeletedAt() int   { return int(u.UpdatedAt.Unix()) }
func (u *User) IsPrivateUser() bool { return u.IsPrivate }
func (u *User) GetBio() string      { return u.Bio }
func (u *User) GetAvatar() string   { return u.Avatar }
func (u *User) GetGroups() []model.GroupModel {
	var result []model.GroupModel
	for _, userGroup := range u.Groups {
		result = append(result, &userGroup)
	}
	return result
}
func (u *User) GetFCMToken() string { return u.FCMToken }

func (u *User) GetStatus() int {
	return u.Status
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
		Email:       args.Email,
		FirebaseUID: args.GoogleId,
		Status:      1,
		Role:        args.Role,
		Username:    args.Username,
	}

	result := repo.db.Create(&userObject)
	return &userObject, result.Error
}

func (repo *repoUserPrivate) Update(userID uint, args map[string]interface{}) (model.UserModel, error) {
	userObject := User{}
	repo.db.First(&userObject, userID)
	if userObject.CreatedAt.IsZero() {
		return nil, errors.New("user not found")
	}

	result := repo.db.Model(&userObject).Updates(args)
	if result.Error != nil {
		return nil, result.Error
	}

	return &userObject, nil
}

func (repo *repoUserPrivate) Delete(id uint) error {
	return repo.db.Delete(&User{}, id).Error
}

func (repo *repoUserPrivate) GetByFirebaseUid(googleId string) (model.UserModel, error) {
	userObject := User{}

	result := repo.db.Where("firebase_uid = ?", googleId).First(&userObject)
	if userObject.CreatedAt.IsZero() {
		return nil, errors.New("user not found")
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

	result := repo.db.
		Preload("Groups").
		Joins("LEFT JOIN group_members ON group_members.member_id = users.id AND group_members.deleted_at IS NULL").
		Preload("Groups.CreatedBy").
		Find(&userObject)
	if userObject.CreatedAt.IsZero() {
		return nil, nil
	}

	return &userObject, result.Error
}

func (repo *repoUserPrivate) GetAll(page int, pageSize int) ([]model.UserModel, error) {
	var users []*User
	offset := (page - 1) * pageSize
	result := repo.db.Offset(offset).Limit(pageSize).Find(&users)

	models := make([]model.UserModel, len(users))
	for i, v := range users {
		models[i] = model.UserModel(v)
	}

	return models, result.Error
}

func (repo *repoUserPrivate) GetAllUserCount() (int64, error) {
	var count int64
	result := repo.db.Model(&User{}).Count(&count)
	return count, result.Error
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
	result := repo.db.Where("LOWER(username) LIKE LOWER(@search) AND status = 1", sql.Named("search", searchParam)).Find(&foundUsers)
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

func (repo *repoUserPrivate) GetAllFCMTokens() ([]string, error) {
	var foundUsers []*User
	result := repo.db.Find(&foundUsers)

	if result.Error != nil {
		return nil, result.Error
	}

	tokens := make([]string, len(foundUsers))
	for i, v := range foundUsers {
		tokens[i] = v.GetFCMToken()
	}
	return tokens, nil
}

func (repo *repoUserPrivate) BanUser(userId uint) (model.UserModel, error) {
	userObject := User{}
	userObject.ID = userId

	result := repo.db.Find(&userObject)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}

	userObject.Status = -1
	return nil, repo.db.Save(&userObject).Error
}

func (repo *repoUserPrivate) UnbanUser(userId uint) (model.UserModel, error) {
	userObject := User{}
	userObject.ID = userId

	result := repo.db.Find(&userObject)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}

	userObject.Status = 1
	return nil, repo.db.Save(&userObject).Error
}

func (repo *repoUserPrivate) UpdateByAdmin(userID uint, args model.AdminUpdateUserRequest) (model.UserModel, error) {
	userObject := User{}
	repo.db.First(&userObject, userID)
	if userObject.CreatedAt.IsZero() {
		return nil, errors.New("user not found")
	}

	// print the args
	fmt.Println(args)

	result := repo.db.Model(&userObject).Updates(args)
	if result.Error != nil {
		return nil, result.Error
	}

	return &userObject, nil
}
