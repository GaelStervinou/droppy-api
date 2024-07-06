package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/http/response_models"
	"go-api/internal/repositories"
	groupservice "go-api/internal/services/group"
	"go-api/pkg/model"
	"net/http"
)

// CreateGroup godoc
//
//	@Summary		Create group
//	@Description	Create group
//	@Tags			group
//	@Accept			json
//	@Produce		json
//	@Param			group	body		model.GroupCreationParam	true	"Group creation object"
//	@Success		201	{object} response_models.GetGroupResponse
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		500
//	@Router			/groups [post]
func CreateGroup(c *gin.Context) {
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

	var groupToCreate model.GroupCreationParam

	if err := c.ShouldBindJSON(&groupToCreate); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	gs := &groupservice.GroupService{
		Repo: repositories.Setup(),
	}

	createdGroup, err := gs.CreateGroup(uintCurrentUserId, groupToCreate)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == createdGroup {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Group not created"})
		return
	}

	user, err := gs.Repo.UserRepository.GetById(uintCurrentUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponse := response_models.FormatGetUserResponse(user)

	groupResponse := response_models.FormatGetGroupResponse(createdGroup, userResponse)

	c.JSON(http.StatusCreated, groupResponse)
}
