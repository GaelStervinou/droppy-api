package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/oov/gothic"
	"go-api/internal/repositories"
	"go-api/internal/services/account"
	"go-api/pkg/jwt_helper"
	"go-api/pkg/model"
	"net/http"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken godoc
//
//	@Summary		Refresh auth token
//	@Description	get a new jwt token from a refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			token body		RefreshTokenRequest	true	"Refresh object"
//	@Success		200	{object} account.TokenInfo
//	@Failure		401
//	@Failure		500
//	@Router			/auth/refresh [get]
func RefreshToken(c *gin.Context) {
	acc := &account.AccountService{
		Repo: repositories.Setup(),
	}

	var refreshToken RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if "" == refreshToken.RefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
		return
	}

	token, err := jwt_helper.VerifyToken(refreshToken.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenInfo, err := acc.LoginFromRefreshToken(token.Raw)
	if err != nil {
		if err.Error() == "refresh token not found" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenInfo)
}

func GoogleAuth(c *gin.Context) {
	err := gothic.BeginAuth(c.Param("provider"), c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func GoogleAuthCallback(c *gin.Context) {
	acc := &account.AccountService{
		Repo: repositories.Setup(),
	}

	user, err := gothic.CompleteAuth(c.Param("provider"), c.Writer, c.Request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if "" == user.Email {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "email is empty"})
		return
	}

	isKnown, err := acc.EmailExists(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if false == isKnown {
		err := acc.CreateWithGoogle(user.Email, user.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tokenInfo, err := acc.LoginWithGoogle(user.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenInfo)
}

// Login godoc
//
//	@Summary		Login
//	@Description	login with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			login body		model.LoginParam	true	"Login object"
//	@Success		200	{object} account.TokenInfo
//	@Failure		422 "Invalid email or password"
//	@Failure		500
//	@Router			/auth [post]
func Login(c *gin.Context) {
	acc := &account.AccountService{
		Repo: repositories.Setup(),
	}
	var loginParam model.LoginParam

	if err := c.ShouldBindJSON(&loginParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenInfo, err := acc.Login(loginParam.Email, loginParam.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, tokenInfo)
}

type FirebaseToken struct {
	IDToken string `json:"id_token"`
}

// FirebaseLogin godoc
//
//	@Summary		Login
//	@Description	login with firebase id token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			idToken formData string true "Firebase ID Token"
//	@Success		200	{object} account.TokenInfo
//	@Failure		422
//	@Failure		500
//	@Router			/auth/oauth_token [post]
func FirebaseLogin(c *gin.Context) {
	token := FirebaseToken{}

	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if "" == token.IDToken {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "id_token is empty"})
		return
	}

	acc := &account.AccountService{
		Repo: repositories.Setup(),
	}

	tokenInfo, err := acc.LoginWithFirebase(token.IDToken, c)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenInfo)
}
