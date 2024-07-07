package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/http/response_models"
	"go-api/internal/repositories"
	groupservice "go-api/internal/services/group"
	"go-api/internal/storage/postgres"
	"go-api/internal/storage/postgres/group"
	"go-api/internal/storage/postgres/user"
	"go-api/pkg/converters"
	"go-api/pkg/errors2"
	"go-api/pkg/model"
	"net/http"
	"strings"
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

	patchedGroup, err := gs.PatchGroup(groupId, uintCurrentUserId, groupPatch)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == patchedGroup {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Group not found"})
		return
	}

	groupResponse := response_models.FormatGetGroupResponse(patchedGroup)

	c.JSON(http.StatusOK, groupResponse)
}

// SearchGroups godoc
//
// @Summary		Search groups
// @Description	Search groups
// @Tags			group
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			search query string true "Search query"
// @Success		200	{object} []response_models.GetSearchGroupResponse
// @Failure		400
// @Failure		500
// @Router			/groups/search [get]
func SearchGroups(c *gin.Context) {
	sqlDB, err := postgres.Connect()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	gr := group.NewRepo(sqlDB)

	query := strings.TrimSpace(c.Query("search"))

	if "" == query {
		c.JSON(400, errors2.MultiFieldsError{Fields: map[string]string{"search": "Search query is required"}})
		return
	}

	searchedGroups, err := gr.Search(query)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if nil == searchedGroups {
		c.JSON(404, gin.H{"error": "Users not found"})
		return
	}

	currentUserId, exists := c.Get("userId")

	var currentUser model.UserModel

	if exists {
		uintCurrentUserId, ok := currentUserId.(uint)
		if ok {
			us := user.NewRepo(sqlDB)
			targetedUser, err := us.GetById(uintCurrentUserId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			currentUser = targetedUser
		}
	}

	gms := &groupservice.GroupMemberService{
		Repo: repositories.Setup(),
	}

	var currentUserGroups []model.GroupModel
	if currentUser != nil {
		currentUserGroups, err = gms.FindAllUserGroups(currentUser.GetID())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	var groupResponse []response_models.GetSearchGroupResponse
	for _, searchedGroup := range searchedGroups {
		isMember := false
		for _, currentUserGroup := range currentUserGroups {
			if searchedGroup.GetID() == currentUserGroup.GetID() {
				isMember = true
				continue
			}
		}
		groupResponse = append(groupResponse, response_models.FormatGetSearchGroupResponse(searchedGroup, isMember))
	}

	c.JSON(200, groupResponse)
}
