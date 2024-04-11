package token

import (
	"context"
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type AuthToken struct {
	gorm.Model
	Token  string `json:"token,omitempty"`
	UserID uint   `json:"userId,omitempty"`
	Expiry int    `json:"expiry,omitempty"`
}

func (s *AuthToken) GetID() uint { return s.ID }

func (s *AuthToken) GetToken() string {
	return s.Token
}

func (s *AuthToken) GetUserID() uint {
	return s.UserID
}

func (s *AuthToken) GetExpiry() int {
	return s.Expiry
}

// Safe checker to know if this file already implements the model interface correctly or not
var _ model.AuthTokenModel = (*AuthToken)(nil)

type repoPrivate struct {
	db *gorm.DB
}

// Safe checker to know if this file already implements the interface correctly or not
var _ model.AuthTokenRepository = (*repoPrivate)(nil)

func NewRepo(db *gorm.DB) model.AuthTokenRepository {
	return &repoPrivate{
		db: db,
	}
}

func (repo *repoPrivate) Create(ctx context.Context, args model.TokenCreationParam) (model.AuthTokenModel, error) {
	tokenObject := AuthToken{
		Token:  args.Token,
		UserID: args.UserID,
		Expiry: args.Expiry,
	}

	result := repo.db.Create(&tokenObject)

	return &tokenObject, result.Error
}

func (repo *repoPrivate) FindByRefreshToken(token string) (model.AuthTokenModel, error) {
	tokenObject := AuthToken{}
	result := repo.db.Where("token = ?", token).First(&tokenObject)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tokenObject, result.Error
}

func (repo *repoPrivate) FindByUserId(userId uint) (model.AuthTokenModel, error) {
	tokenObject := AuthToken{}
	result := repo.db.Where("user_id = ?", userId).First(&tokenObject)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tokenObject, result.Error
}

func (repo *repoPrivate) Delete(recordId uint) error {
	return repo.db.Delete(&AuthToken{}, recordId).Error
}
