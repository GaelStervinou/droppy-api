package account

import (
	"context"
	"encoding/base64"
	"errors"
	"go-api/internal/repositories"
	"go-api/pkg/jwt_helper"
	"go-api/pkg/model"
	"go-api/pkg/random"
	"go-api/pkg/services/account"
	"go-api/pkg/validation"
	"golang.org/x/crypto/argon2"
	"time"
)

type AccountService struct {
	Repo *repositories.Repositories
}

func (a *AccountService) Create(firstname string, lastname string, email string, password string, username string) error {
	user := model.UserCreationParam{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Password:  password,
		Username:  username,
		//TODO peut-être passé par une struct pour le role ou au moins un enum ?
		Roles: []string{"user", "admin"},
	}
	validationError := validation.ValidateUserCreation(user)
	if len(validationError.Fields) > 0 {
		return validationError
	}
	hashedPassword := argonFromPassword(password)

	user.Password = hashedPassword

	_, err := a.Repo.UserRepository.Create(user)

	return err
}

func (a *AccountService) CreateWithGoogle(firstname string, lastname string, email string, googleId string) error {
	_, err := a.Repo.UserRepository.CreateWithGoogle(
		model.UserCreationWithGoogleParam{
			Firstname: firstname,
			Lastname:  lastname,
			Email:     email,
			GoogleId:  googleId,
			Username:  random.RandStringRunes(10),
			//TODO peut-être passé par une struct pour le role ou au moins un enum ?
			Roles: []string{"user"},
		},
	)

	return err
}

func (a *AccountService) Login(email string, password string) (*account.TokenInfo, error) {
	user, err := a.Repo.UserRepository.GetByEmail(email)
	if err != nil {
		return &account.TokenInfo{}, err
	}

	if user.GetPassword() != argonFromPassword(password) {
		return &account.TokenInfo{}, errors.New("email or password does not match our record")
	}

	newToken, refreshToken, newTokenExpiry, err := jwt_helper.GenerateToken(user.GetID(), user.GetRoles())
	if err != nil {
		return &account.TokenInfo{}, err
	}

	_, err = a.Repo.TokenRepository.Create(context.TODO(), model.TokenCreationParam{
		Token:  refreshToken,
		UserID: user.GetID(),
		Expiry: newTokenExpiry,
	})

	if err != nil {
		return &account.TokenInfo{}, err
	}

	return &account.TokenInfo{JWTToken: newToken, RefreshToken: refreshToken, Expiry: newTokenExpiry}, nil
}

func (a *AccountService) LoginFromRefreshToken(refreshToken string) (*account.TokenInfo, error) {
	t, err := a.Repo.TokenRepository.FindByRefreshToken(refreshToken)
	if err != nil {
		return &account.TokenInfo{}, err
	}
	if t == nil {
		return &account.TokenInfo{}, errors.New("refresh token not found")
	}

	if t.GetExpiry() < int(time.Now().Unix()) {
		return &account.TokenInfo{}, errors.New("refresh token expired")
	}

	user, err := a.Repo.UserRepository.GetById(t.GetUserID())
	if err != nil {
		return &account.TokenInfo{}, err
	}

	newToken, newRefreshToken, newTokenExpiry, err := jwt_helper.GenerateToken(user.GetID(), user.GetRoles())
	if err != nil {
		return &account.TokenInfo{}, err
	}

	_, err = a.Repo.TokenRepository.Create(context.TODO(), model.TokenCreationParam{
		Token:  newRefreshToken,
		UserID: user.GetID(),
		Expiry: newTokenExpiry,
	})

	if err != nil {
		return &account.TokenInfo{}, err
	}

	err = a.Repo.TokenRepository.Delete(t.GetID())

	return &account.TokenInfo{JWTToken: newToken, RefreshToken: newRefreshToken, Expiry: newTokenExpiry}, nil
}

func (a *AccountService) LoginWithGoogle(email string) (*account.TokenInfo, error) {
	user, err := a.Repo.UserRepository.GetByEmail(email)
	if err != nil {
		return &account.TokenInfo{}, err
	}

	//TODO refacto avec function login juste au dessus
	newToken, refreshToken, newTokenExpiry, err := jwt_helper.GenerateToken(user.GetID(), user.GetRoles())
	if err != nil {
		return &account.TokenInfo{}, err
	}

	_, err = a.Repo.TokenRepository.Create(context.TODO(), model.TokenCreationParam{
		Token:  refreshToken,
		UserID: user.GetID(),
		Expiry: newTokenExpiry,
	})

	if err != nil {
		return &account.TokenInfo{}, err
	}

	return &account.TokenInfo{JWTToken: newToken, RefreshToken: refreshToken, Expiry: newTokenExpiry}, nil
}

func (a *AccountService) Logout(userId uint) error {
	return nil
}

func (a *AccountService) EmailExists(email string) (bool, error) {
	_, err := a.Repo.UserRepository.GetByEmail(email)
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
	//TODO un salt par user ? ou un salt global ?
	salt := []byte("salt1234")

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return base64.RawStdEncoding.EncodeToString(hash)
}
