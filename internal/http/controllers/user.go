package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-api/internal/repositories"
	"go-api/internal/services/account"
	"go-api/internal/storage/postgres"
	"go-api/internal/storage/postgres/user"
	"go-api/pkg/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetUserById godoc
//
// @Summary		Get user by ID
// @Description	Get user by ID
// @Tags			user
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			id path int true "User ID"
// @Success		200	{object} user.User
// @Failure		400
// @Failure		403
// @Failure		404
// @Failure		500
// @Router			/users/{id} [get]
func GetUserById(c *gin.Context) {
	sqlDB, err := postgres.Connect()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	us := user.NewRepo(sqlDB)

	id := strings.TrimSpace(c.Param("id"))

	if "" == id {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	requestedUser, err := us.GetById(uint(userID))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if nil == requestedUser {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	currentUserId, exists := c.Get("userId")

	uintCurrentUserId := uint(0)
	if exists {
		toUint, ok := currentUserId.(uint)
		fmt.Println(toUint, ok)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		uintCurrentUserId = toUint
	}

	email := requestedUser.GetEmail()
	emailPointer := &email
	phoneNumber := requestedUser.GetPhoneNumber()
	phoneNumberPointer := &phoneNumber
	bio := requestedUser.GetBio()
	bioPointer := &bio
	if "" == bio {
		bioPointer = nil
	}
	avatar := requestedUser.GetAvatar()
	avatarPointer := &avatar
	if "" == avatar {
		avatarPointer = nil
	}
	createdAt := time.Unix(int64(requestedUser.GetCreatedAt()), 0)
	updatedAt := time.Unix(int64(requestedUser.GetUpdatedAt()), 0)
	createdAtPointer := &createdAt
	updatedAtPointer := &updatedAt

	userResponse := UserResponse{
		ID:          requestedUser.GetID(),
		GoogleID:    requestedUser.GetGoogleID(),
		Email:       emailPointer,
		Username:    requestedUser.GetUsername(),
		Firstname:   requestedUser.GetFirstname(),
		Lastname:    requestedUser.GetLastname(),
		PhoneNumber: phoneNumberPointer,
		Bio:         bioPointer,
		Avatar:      avatarPointer,
		IsPrivate:   requestedUser.IsPrivateUser(),
		Role:        requestedUser.GetRole(),
		CreatedAt:   createdAtPointer,
		UpdatedAt:   updatedAtPointer,
	}

	fmt.Println(uintCurrentUserId, userResponse.ID)
	if uintCurrentUserId != userResponse.ID {
		userResponse.HidePersonalInfo()
		fmt.Println(userResponse)
	}

	c.JSON(200, userResponse)
}

// Create godoc
//
//	@Summary		Create user
//	@Description	Create user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model.UserCreationParam	true	"User creation object"
//	@Success		201	{object} account.TokenInfo
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		500
//	@Router			/users [post]
func Create(c *gin.Context) {
	sqlDB, err := postgres.Connect()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var userToCreate model.UserCreationParam

	if err := c.ShouldBindJSON(&userToCreate); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	userToCreate.Role = "user"
	us := user.NewRepo(sqlDB)
	createdUser, err := us.Create(userToCreate)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == createdUser {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not created"})
		return
	}

	acc := &account.AccountService{
		Repo: repositories.Setup(),
	}

	tokenInfo, err := acc.Login(createdUser.GetEmail(), userToCreate.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tokenInfo)
}

// PatchUserById godoc
//
// @Summary		Patch user by ID
// @Description	Patch user by ID
// @Tags			user
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			id path int true "User ID"
//
//	@Param			user	body		model.UserPatchParam	true	"User creation object"
//
// @Success		200	{object} user.User
// @Failure		400
// @Failure		403
// @Failure		404
// @Failure		500
// @Router			/users/{id} [patch]
func PatchUserById(c *gin.Context) {
	currentUserId, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	sqlDB, err := postgres.Connect()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	us := user.NewRepo(sqlDB)

	id := strings.TrimSpace(c.Param("id"))

	if "" == id {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	if currentUserId != uint(userID) {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	requestedUser, err := us.GetById(uint(userID))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if nil == requestedUser {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var userToPatch model.UserPatchParam

	if err := c.ShouldBindJSON(&userToPatch); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	userToPatch.Email = requestedUser.GetEmail()

	updatedUser, err := us.Update(userToPatch)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == updatedUser {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not updated"})
		return
	}

	c.JSON(200, requestedUser)
}

type UserResponse struct {
	ID          uint
	GoogleID    *string
	Email       *string
	Username    string
	Firstname   string
	Lastname    string
	PhoneNumber *string
	Bio         *string
	Avatar      *string
	IsPrivate   bool
	Role        string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (u *UserResponse) HidePersonalInfo() {
	u.Email = nil
	u.PhoneNumber = nil
	u.GoogleID = nil
	u.CreatedAt = nil
	u.UpdatedAt = nil
}
