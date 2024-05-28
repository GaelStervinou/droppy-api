package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/oov/gothic"
	"go-api/internal/repositories"
	"go-api/internal/services/account"
	"go-api/pkg/jwt_helper"
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
		err := acc.CreateWithGoogle(user.FirstName, user.LastName, user.Email, user.UserID)
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
//	@Param			email formData string true "Email"
//	@Param			password formData string true "Password"
//	@Success		200	{object} account.TokenInfo
//	@Failure		422 "Invalid email or password"
//	@Failure		500
//	@Router			/auth/login [post]
func Login(c *gin.Context) {
	acc := &account.AccountService{
		Repo: repositories.Setup(),
	}

	email := c.PostForm("email")
	password := c.PostForm("password")

	tokenInfo, err := acc.Login(email, password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, tokenInfo)
}
