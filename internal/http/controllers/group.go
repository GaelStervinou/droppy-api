package controllers

import (
	"errors"
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

// JoinGroup godoc
//
//	@Summary		Join Group
//	@Description	Join Group
//	@Tags			group
//	@Accept			json
//
// @Security BearerAuth
//
//	@Produce		json
//	@Param			group	body		model.GroupMemberCreationParam	true	"Join group creation object"
//	@Success		201	{object} response_models.GetGroupMemberResponse
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		500
//	@Router			/groups/members/{id}/join [post]
func JoinGroup(c *gin.Context) {
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

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var groupMemberToCreate model.GroupMemberCreationParam

	if err := c.ShouldBindJSON(&groupMemberToCreate); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	groupMemberToCreate.GroupID = groupId

	gms := &groupservice.GroupMemberService{
		Repo: repositories.Setup(),
	}

	createdGroupMember, err := gms.JoinGroup(uintCurrentUserId, uintCurrentUserId, groupMemberToCreate)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == createdGroupMember {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Group member not created"})
		return
	}

	groupMemberResponse := response_models.FormatGetGroupMemberResponse(createdGroupMember)

	c.JSON(http.StatusCreated, groupMemberResponse)
}

// DeleteGroupMember godoc
//
//	@Summary		Delete Group Member
//	@Description	Delete Group Member
//	@Tags			group
//	@Accept			json
//
// @Security BearerAuth
//
//	@Produce		json
//	@Success		204
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		400
//	@Failure		500
//	@Router			/groups/members/{groupId}/{memberId} [delete]
func DeleteGroupMember(c *gin.Context) {
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

	groupId := c.Param("groupId")
	memberId := c.Param("memberId")

	if "" == groupId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "groupId is required"})
		return
	}

	if "" == memberId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "memberId is required"})
		return
	}

	groupIdUint, err := converters.StringToUint(groupId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	memberIdUint, err := converters.StringToUint(memberId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	gms := &groupservice.GroupMemberService{
		Repo: repositories.Setup(),
	}

	err = gms.DeleteGroupMember(uintCurrentUserId, groupIdUint, memberIdUint)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Group member deleted"})
}

// PatchGroupMember godoc
//
//	@Summary		Patch Group Member
//	@Description	Patch Group Member
//	@Tags			group
//	@Accept			json
//
// @Security BearerAuth
//
//	@Produce		json
//	@Param			group	body		model.GroupMemberPatchParam	true	"Group member patch object"
//	@Success		200	{object} response_models.GetGroupMemberResponse
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		400
//	@Failure		500
//	@Router			/groups/members/{groupId}/{memberId} [patch]
func PatchGroupMember(c *gin.Context) {
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

	groupId := c.Param("groupId")
	memberId := c.Param("memberId")

	if "" == groupId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "groupId is required"})
		return
	}

	if "" == memberId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "memberId is required"})
		return
	}

	groupIdUint, err := converters.StringToUint(groupId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	memberIdUint, err := converters.StringToUint(memberId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	var groupMemberPatch model.GroupMemberPatchParam

	if err := c.ShouldBindJSON(&groupMemberPatch); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	gms := &groupservice.GroupMemberService{
		Repo: repositories.Setup(),
	}

	patchedGroupMember, err := gms.UpdateGroupMemberRole(uintCurrentUserId, groupIdUint, memberIdUint, groupMemberPatch)

	if err != nil {
		var notAllowedErr errors2.NotAllowedError
		if errors.As(err, &notAllowedErr) {
			c.JSON(http.StatusForbidden, gin.H{"error": notAllowedErr.Reason})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == patchedGroupMember {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Group member not found"})
		return
	}

	groupMemberResponse := response_models.FormatGetGroupMemberResponse(patchedGroupMember)

	c.JSON(http.StatusOK, groupMemberResponse)
}

// AcceptGroupMemberRequest godoc
//
//	@Summary		Accept Group Member Request
//	@Description	Accept Group Member Request
//	@Tags			group
//	@Accept			json
//
// @Security BearerAuth
//
//	@Produce		json
//	@Success		200	{object} response_models.GetGroupMemberResponse
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		400
//	@Failure		500
//	@Router			/group/{id}/{memberId}/accept [post]
func AcceptGroupMemberRequest(c *gin.Context) {
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

	groupId := c.Param("id")
	memberId := c.Param("memberId")

	if "" == groupId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "groupId is required"})
		return
	}

	if "" == memberId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "memberId is required"})
		return
	}

	groupIdUint, err := converters.StringToUint(groupId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	memberIdUint, err := converters.StringToUint(memberId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	gms := &groupservice.GroupMemberService{
		Repo: repositories.Setup(),
	}

	groupMember, err := gms.AcceptGroupMember(uintCurrentUserId, groupIdUint, memberIdUint)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == groupMember {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Group member not found"})
		return
	}

	groupMemberResponse := response_models.FormatGetGroupMemberResponse(groupMember)

	c.JSON(http.StatusOK, groupMemberResponse)
}

// RefuseGroupMemberRequest godoc
//
//	@Summary		Refuse Group Member Request
//	@Description	Refuse Group Member Request
//	@Tags			group
//	@Accept			json
//
// @Security BearerAuth
//
//	@Produce		json
//	@Success		204	{} No Content
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		400
//	@Failure		500
//	@Router			/groups/members/{id}/{memberId}/refuse/ [post]
func RefuseGroupMemberRequest(c *gin.Context) {
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

	groupId := c.Param("id")
	memberId := c.Param("memberId")

	if "" == groupId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "groupId is required"})
		return
	}

	if "" == memberId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "memberId is required"})
		return
	}

	groupIdUint, err := converters.StringToUint(groupId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	memberIdUint, err := converters.StringToUint(memberId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	gms := &groupservice.GroupMemberService{
		Repo: repositories.Setup(),
	}

	err = gms.DeleteGroupMember(uintCurrentUserId, groupIdUint, memberIdUint)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Group member deleted"})
}

// GetPendingGroupMemberRequests godoc
//
//	@Summary		Get Pending Group Member Requests
//	@Description	Get Pending Group Member Requests
//	@Tags			group
//	@Accept			json
//
// @Security BearerAuth
//
//	@Produce		json
//	@Success		200	{object} []response_models.GetGroupMemberResponse
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		400
//	@Failure		500
//	@Router			/groups/members/{id}/pending [get]
func GetPendingGroupMemberRequests(c *gin.Context) {
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

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	gms := &groupservice.GroupMemberService{
		Repo: repositories.Setup(),
	}

	groupMembers, err := gms.GetPendingGroupMemberRequests(uintCurrentUserId, groupId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var groupMemberResponse []response_models.GetGroupMemberResponse
	for _, groupMember := range groupMembers {
		groupMemberResponse = append(groupMemberResponse, response_models.FormatGetGroupMemberResponse(groupMember))
	}

	c.JSON(http.StatusOK, groupMemberResponse)
}
