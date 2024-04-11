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
	Find(ctx context.Context, token string) (AuthTokenModel, error)
	Delete(ctx context.Context, userId uint) error
}
