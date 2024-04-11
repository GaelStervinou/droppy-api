package model

import "context"

type TokenCreationParam struct {
	Token  string
	UserID uint
	Expiry int
}

type AuthTokenModel interface {
	GetID() uint
	GetToken() string
	GetUserID() uint
	GetExpiry() int
}

type AuthTokenRepository interface {
	Create(ctx context.Context, args TokenCreationParam) (AuthTokenModel, error)
	FindByRefreshToken(token string) (AuthTokenModel, error)
	FindByUserId(uint) (AuthTokenModel, error)
	Delete(uint) error
}
