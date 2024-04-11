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
	existingRow, err := repo.Find(ctx, args.Token)
	if err != nil {
		return nil, err
	}

	var result *gorm.DB
	if existingRow != nil {
		result = repo.db.Create(&tokenObject)
	} else {
		result = repo.db.Save(&tokenObject)
	}

	return &tokenObject, result.Error
}

func (repo *repoPrivate) Find(ctx context.Context, token string) (model.AuthTokenModel, error) {
	tokenObject := AuthToken{Token: token}
	result := repo.db.Find(&tokenObject)

	return &tokenObject, result.Error
}

func (repo *repoPrivate) Delete(ctx context.Context, userId uint) error {
	return repo.db.Delete(&AuthToken{}, userId).Error
}
