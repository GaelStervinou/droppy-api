package account

import (
	"context"
	"encoding/base64"
	"errors"
	"go-api/internal/repositories"
	"go-api/pkg/jwt_helper"
	"go-api/pkg/model"
	"go-api/pkg/services/account"
	"golang.org/x/crypto/argon2"
)

type AccountService struct {
	Repo *repositories.Repositories
}

func (a *AccountService) Create(ctx context.Context, firstname string, lastname string, email string, password string) error {
	hashedPassword := argonFromPassword(password)

	_, err := a.Repo.UserRepository.Create(
		context.TODO(),
		model.UserCreationParam{
			Firstname: firstname,
			Lastname:  lastname,
			Email:     email,
			Password:  hashedPassword,
			//TODO peut-être passé par une struct pour le role ou au moins un enum ?
			Role: "user",
		},
	)

	return err
}

func (a *AccountService) CreateWithGoogle(ctx context.Context, firstname string, lastname string, email string, googleId string) error {
	_, err := a.Repo.UserRepository.CreateWithGoogle(
		context.TODO(),
		model.UserCreationWithGoogleParam{
			Firstname: firstname,
			Lastname:  lastname,
			Email:     email,
			GoogleId:  googleId,
			//TODO peut-être passé par une struct pour le role ou au moins un enum ?
			Role: "user",
		},
	)

	return err
}

func (a *AccountService) Login(ctx context.Context, email string, password string) (*account.TokenInfo, error) {
	//first read the user data
	user, err := a.Repo.UserRepository.GetByEmail(context.TODO(), email)
	if err != nil {
		return &account.TokenInfo{}, err
	}

	//then, compare password
	if user.GetPassword() != argonFromPassword(password) {
		return &account.TokenInfo{}, errors.New("email or password does not match our record")
	}

	//login successful, now generate a random token
	newToken, newTokenExpiry, err := jwt_helper.GenerateToken(user.GetID())
	if err != nil {
		return &account.TokenInfo{}, err
	}

	_, err = a.Repo.TokenRepository.Create(context.TODO(), model.TokenCreationParam{
		Token:  newToken,
		UserID: user.GetID(),
		Email:  email,
		Expiry: newTokenExpiry,
	})

	if err != nil {
		return &account.TokenInfo{}, err
	}

	return &account.TokenInfo{Token: newToken, Expiry: newTokenExpiry}, nil
}

func (a *AccountService) LoginWithGoogle(ctx context.Context, email string) (*account.TokenInfo, error) {
	user, err := a.Repo.UserRepository.GetByEmail(context.TODO(), email)
	if err != nil {
		return &account.TokenInfo{}, err
	}

	//TODO refacto avec function login juste au dessus
	newToken, newTokenExpiry, err := jwt_helper.GenerateToken(user.GetID())
	if err != nil {
		return &account.TokenInfo{}, err
	}

	_, err = a.Repo.TokenRepository.Create(context.TODO(), model.TokenCreationParam{
		Token:  newToken,
		UserID: user.GetID(),
		Email:  email,
		Expiry: newTokenExpiry,
	})

	if err != nil {
		return &account.TokenInfo{}, err
	}

	return &account.TokenInfo{Token: newToken, Expiry: newTokenExpiry}, nil
}

func (a *AccountService) Logout(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AccountService) EmailExists(ctx context.Context, email string) (bool, error) {
	_, err := a.Repo.UserRepository.GetByEmail(context.TODO(), email)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// Safe checker to know if this file already implements the interface correctly or not
var _ account.AccountServiceIface = (*AccountService)(nil)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func argonFromPassword(password string) string {
	p := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  8,
		keyLength:   16,
	}
	//TODO utiliser var d'env pour le salt ?
	salt := []byte("salt1234")

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return base64.RawStdEncoding.EncodeToString(hash)
}
