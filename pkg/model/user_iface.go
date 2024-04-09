package model

import (
	"context"
)

type UserModel interface {
	GetID() uint
	GetGoogleID() string
	GetEmail() string
	GetPassword() string
	GetFirstname() string
	GetLastname() string
	GetRole() string
	GetCreatedAt() int
	GetUpdatedAt() int
}

type UserRepository interface {
	GetByGoogleAuthId(ctx context.Context, googleID string) (UserModel, error)
	GetByEmail(ctx context.Context, email string) (UserModel, error)
	GetById(ctx context.Context, id uint) (UserModel, error)
	Create(ctx context.Context, args UserCreationParam) (UserModel, error)
	CreateWithGoogle(ctx context.Context, args UserCreationWithGoogleParam) (UserModel, error)
	Update(ctx context.Context, user UserModel) (UserModel, error)
	Delete(ctx context.Context, id uint) error
	GetAll() ([]UserModel, error)
}

type UserCreationParam struct {
	Firstname string
	Lastname  string
	Email     string
	Password  string
	Role      string
}

type UserCreationWithGoogleParam struct {
	Firstname string
	Lastname  string
	Email     string
	GoogleId  string
	Role      string
}
