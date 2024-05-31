package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/storage/postgres"
	"go-api/internal/storage/postgres/follow"
	"go-api/internal/storage/postgres/user"
	"go-api/pkg/converters"
	"go-api/pkg/model"
	"net/http"
)

// FollowUser godoc
//
// @Summary		Follow user
// @Description	Follow user by its ID
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security BearerAuth
//
//	@Param			user	body		model.FollowCreationParam	true	"Follow creation object"
//
// @Success		201	{object} follow.Follow
// @Failure		422
// @Failure		401
// @Failure		500
// @Router			/follows [post]
func FollowUser(c *gin.Context) {
	var followCreationParam model.FollowCreationParam

	if err := c.ShouldBindJSON(&followCreationParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	sqlDB, err := postgres.Connect()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	us := user.NewRepo(sqlDB)

	requestedUser, err := us.GetById(followCreationParam.UserToFollowID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == requestedUser {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "User not found"})
		return
	}

	if uintCurrentUserId == followCreationParam.UserToFollowID {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "You can't follow yourself"})
		return
	}

	isFollowingAllowed, err := us.CanUserBeFollowed(followCreationParam.UserToFollowID)

	if err != nil || !isFollowingAllowed {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "You can't follow this user for now. Please try again later"})
		return
	}

	followRepo := follow.NewRepo(sqlDB)

	if alreadyFollowing, _ := followRepo.AreAlreadyFollowing(uintCurrentUserId, followCreationParam.UserToFollowID); alreadyFollowing {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "You are already following this user"})
		return
	}

	createdFollow, err := followRepo.Create(uintCurrentUserId, followCreationParam.UserToFollowID, !requestedUser.IsPrivateUser())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if nil == createdFollow {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not followed"})
		return
	}

	c.JSON(http.StatusCreated, createdFollow)
}

// GetMyFollowers godoc
//
// @Summary		 Get my followers
// @Description	 Get all followers of the current user
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security BearerAuth
//
// @Success		201	{object} []user.User
// @Failure		422
// @Failure		401
// @Failure		500
// @Router			/follows/my-followers [get]
func GetMyFollowers(c *gin.Context) {
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

	sqlDB, err := postgres.Connect()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	followRepo := follow.NewRepo(sqlDB)

	followers, err := followRepo.GetFollowers(uintCurrentUserId)

	followerIds := make([]uint, 0)
	for _, follower := range followers {
		followerIds = append(followerIds, follower.GetFollowerID())
	}


	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userRepo := user.NewRepo(sqlDB)

	users, err := userRepo.GetUsersFromUserIds(followerIds)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if nil == users {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetMyPendingRequests godoc
//
// @Summary		 Get my pending requests
// @Description	 Get all pending requests of the current user
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security BearerAuth
//
// @Success		201	{object} []user.User
// @Failure		422
// @Failure		401
// @Failure		500
// @Router			/follows/pending [get]
func GetMyPendingRequests(c *gin.Context) {
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

	sqlDB, err := postgres.Connect()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	followRepo := follow.NewRepo(sqlDB)

	pendingRequests, err := followRepo.GetPendingRequests(uintCurrentUserId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pendingRequests)
}

// AcceptRequest godoc
//
// @Summary		 Accept request
// @Description	 Accept request by its ID
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security BearerAuth
//
// @Success		201	{object} []user.User
// @Failure		422
// @Failure		401
// @Failure		500
// @Router			/follows/accept/{id} [post]
func AcceptRequest(c *gin.Context) {
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

	sqlDB, err := postgres.Connect()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	if "" == id {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	followId, err := converters.StringToUint(id)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid follow ID"})
		return
	}

	followRepo := follow.NewRepo(sqlDB)

	IsMyFollow, err := followRepo.IsMyFollow(uintCurrentUserId, followId)

	if IsMyFollow == false {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	err = followRepo.AcceptRequest(followId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Follow request accepted"})
}

// RejectRequest godoc
//
// @Summary		 Refuse follow request
// @Description	 Refuse follow request by its ID
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security BearerAuth
//
// @Success		201	{object} []user.User
// @Failure		422
// @Failure		401
// @Failure		500
// @Router			/follows/refuse/{id} [post]
func RejectRequest(c *gin.Context) {
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

	sqlDB, err := postgres.Connect()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	if "" == id {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	followId, err := converters.StringToUint(id)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid follow ID"})
		return
	}

	followRepo := follow.NewRepo(sqlDB)

	IsMyFollow, err := followRepo.IsMyFollow(uintCurrentUserId, followId)

	if IsMyFollow == false {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	err = followRepo.RejectRequest(followId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Follow request refused"})
}