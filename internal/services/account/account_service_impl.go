package account

import (
	"context"
	"errors"
	"go-api/internal/repositories"
	"go-api/internal/storage/firebase"
	"go-api/pkg/hash"
	"go-api/pkg/jwt_helper"
	"go-api/pkg/model"
	"go-api/pkg/random"
	"go-api/pkg/services/account"
	"go-api/pkg/validation"
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
		Role:      "user",
	}
	validationError := validation.ValidateUserCreation(user)
	if len(validationError.Fields) > 0 {
		return validationError
	}
	hashedPassword, err := hash.GenerateFromPassword(password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	_, err = a.Repo.UserRepository.Create(user)

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
			Role:      "user",
		},
	)
	return err
}

func (a *AccountService) Login(email string, password string) (*account.TokenInfo, error) {
	user, err := a.Repo.UserRepository.GetByEmail(email)
	if err != nil {
		return &account.TokenInfo{}, err
	}
	match, err := hash.ComparePasswordAndHash(password, user.GetPassword())
	if err != nil {
		return &account.TokenInfo{}, errors.New("Error while comparing password and hash")
	}
	if !match {
		return &account.TokenInfo{}, errors.New("email or password does not match our record")
	}

	newToken, refreshToken, newTokenExpiry, err := jwt_helper.GenerateToken(user.GetID(), user.GetRole())
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

func (a *AccountService) LoginWithFirebase(token string, ctx context.Context) (*account.TokenInfo, error) {
	firebaseRepo, err := firebase.NewRepo()

	if err != nil {
		return &account.TokenInfo{}, err
	}

	client, err := firebaseRepo.App.Auth(ctx)
	if err != nil {
		return &account.TokenInfo{}, err
	}

	decodedToken, err := client.VerifyIDToken(ctx, token)

	if err != nil {
		return &account.TokenInfo{}, err
	}

	user, err := a.Repo.UserRepository.GetByGoogleAuthId(decodedToken.UID)

	if err != nil {
		if "user not found" == err.Error() {
			err = a.CreateWithGoogle(decodedToken.Claims["name"].(string), decodedToken.Claims["name"].(string), decodedToken.Claims["email"].(string), decodedToken.UID)
			if err != nil {
				return &account.TokenInfo{}, err
			}

			user, err = a.Repo.UserRepository.GetByGoogleAuthId(decodedToken.UID)

			if err != nil {
				return &account.TokenInfo{}, err
			}
		} else {
			return &account.TokenInfo{}, err
		}
	}

	return a.LoginWithGoogle(user.GetEmail())
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

	newToken, newRefreshToken, newTokenExpiry, err := jwt_helper.GenerateToken(user.GetID(), user.GetRole())
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
	newToken, refreshToken, newTokenExpiry, err := jwt_helper.GenerateToken(user.GetID(), user.GetRole())
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
