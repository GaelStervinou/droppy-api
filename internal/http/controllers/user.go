package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-api/internal/http/response_models"
	"go-api/internal/repositories"
	"go-api/internal/services/account"
	"go-api/internal/services/user"
	"go-api/internal/storage/postgres"
	"go-api/pkg/errors2"
	"go-api/pkg/model"
	"net/http"
	"strconv"
	"strings"
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
// @Success		200	{object} response_models.GetUserResponse
// @Failure		400
// @Failure		403
// @Failure		404
// @Router			/users/{id} [get]
func GetUserById(c *gin.Context) {
	currentUserId, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uintCurrentUserId, ok := currentUserId.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sqlDB := postgres.Connect()

	us := postgres.NewUserRepo(sqlDB)

	id := strings.TrimSpace(c.Param("id"))

	if "" == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	requestedUser, err := us.GetById(uint(userID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if nil == requestedUser || requestedUser.GetStatus() != 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	dr := postgres.NewDropRepo(sqlDB)

	pinnedDrops, err := dr.GetUserPinnedDrops(requestedUser.GetID())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dnr := postgres.NewDropNotifRepo(sqlDB)

	lastNotification, err := dnr.GetCurrentDropNotification()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userLastDrop, err := dr.GetUserLastDrop(requestedUser.GetID(), lastNotification.GetID())
	if err != nil {
		if err.Error() != "record not found" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userLastDrop = nil
	}
	isLastDropLiking := false
	if nil != userLastDrop {
		lr := postgres.NewLikeRepo(sqlDB)
		isLastDropLiking, err = lr.LikeExists(userLastDrop.GetID(), requestedUser.GetID())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	fr := postgres.NewFollowRepo(sqlDB)

	totalFollowers := fr.CountFollowers(requestedUser.GetID())
	totalFollowing := fr.CountFollowed(requestedUser.GetID())
	totalDrops := dr.CountUserDrops(requestedUser.GetID())
	currentFollow, err := fr.GetUserFollowedBy(uintCurrentUserId, requestedUser.GetID())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponse := response_models.FormatGetOneUserResponse(
		requestedUser,
		userLastDrop,
		isLastDropLiking,
		pinnedDrops,
		totalFollowers,
		totalFollowing,
		totalDrops,
		currentFollow,
		uintCurrentUserId,
	)

	c.JSON(http.StatusOK, userResponse)
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
	sqlDB := postgres.Connect()

	var userToCreate model.UserCreationParam

	if err := c.ShouldBindJSON(&userToCreate); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	userToCreate.Role = "user"
	us := postgres.NewUserRepo(sqlDB)
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

	tokenInfo, err := acc.Login(createdUser.GetEmail(), userToCreate.Password, "")
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
// @Param			user	body		model.UserPatchParam	true	"User creation object"
//
// @Success		200	{object} postgres.User
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

	us := user.NewUserService(repositories.Setup())

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

	requestedUser, err := us.Repo.UserRepository.GetById(uint(userID))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if nil == requestedUser {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var userToPatch model.UserPatchParam

	if err := c.MustBindWith(&userToPatch, binding.FormMultipart); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := us.UpdateUser(uint(userID), userToPatch)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == updatedUser {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not updated"})
		return
	}

	c.JSON(200, updatedUser)
}

// SearchUsers godoc
//
// @Summary		Search users
// @Description	Search users
// @Tags			user
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			search query string true "Search query"
// @Success		200	{object} []response_models.GetUserResponse
// @Failure		400
// @Failure		500
// @Router			/users/search [get]
func SearchUsers(c *gin.Context) {
	sqlDB := postgres.Connect()

	us := postgres.NewUserRepo(sqlDB)

	query := strings.TrimSpace(c.Query("search"))

	if "" == query {
		c.JSON(400, errors2.MultiFieldsError{Fields: map[string]string{"search": "Search query is required"}})
		return
	}

	users, err := us.Search(query)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if nil == users {
		c.JSON(404, gin.H{"error": "Users not found"})
		return
	}

	var usersResponse []response_models.GetUserResponseInterface
	for _, searchedUser := range users {
		usersResponse = append(usersResponse, response_models.FormatGetUserResponse(searchedUser))
	}

	c.JSON(200, usersResponse)
}
