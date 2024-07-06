package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/http/response_models"
	"go-api/internal/repositories"
	groupservice "go-api/internal/services/group"
	"go-api/pkg/converters"
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
//
// @Security BearerAuth
//
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

	groupResponse := response_models.FormatGetGroupResponse(createdGroup)

	c.JSON(http.StatusCreated, groupResponse)
}

// PatchGroup godoc
//
//	@Summary		Patch group
//	@Description	Patch group
//	@Tags			group
//	@Accept			json
//
// @Security BearerAuth
// @Param			id path int true "User ID"
//
//	@Produce		json
//	@Param			group	body		model.GroupPatchParam	true	"Group patch object"
//	@Success		200	{object} response_models.GetGroupResponse
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		500
//	@Router			/groups/{id} [patch]
func PatchGroup(c *gin.Context) {
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

	id := c.Param("id")

	if "" == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	groupId, err := converters.StringToUint(id)

	var groupPatch model.GroupPatchParam

	if err := c.ShouldBindJSON(&groupPatch); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	gs := &groupservice.GroupService{
		Repo: repositories.Setup(),
	}

	group, err := gs.PatchGroup(groupId, uintCurrentUserId, groupPatch)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == group {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Group not found"})
		return
	}

	groupResponse := response_models.FormatGetGroupResponse(group)

	c.JSON(http.StatusOK, groupResponse)
}
