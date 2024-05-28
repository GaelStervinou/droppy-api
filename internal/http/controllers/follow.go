package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/storage/postgres"
	"go-api/internal/storage/postgres/follow"
	"go-api/internal/storage/postgres/user"
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
