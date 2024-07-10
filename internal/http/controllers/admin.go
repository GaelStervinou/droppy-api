package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-api/internal/http/response_models"
	"go-api/internal/storage/postgres"
)

// GetAllUsers godoc
//
// @Summary		Get all users
// @Description	Get all users by admin user
// @Tags			user
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Success		200	{object} []response_models.GetUserResponse
// @Failure		500
// @Router			/admin/users [get]
func GetAllUsers(c *gin.Context) {
	sqlDB, err := postgres.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	us := postgres.NewUserRepo(sqlDB)

	users, err := us.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}

	var usersResponse []response_models.GetUserResponseInterface
	for _, userModel := range users {
		usersResponse = append(usersResponse, response_models.FormatAdminGetUserResponse(userModel))
	}

	c.JSON(http.StatusOK, usersResponse)
}

// GetAllGroups godoc
//
// @Summary		Get all groups
// @Description	Get all groups by admin user
// @Tags			user
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Success		200	{object} []response_models.GetGroupResponse
// @Failure		500
// @Router			/admin/groups [get]
func GetAllGroups(c *gin.Context) {
	sqlDB, err := postgres.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gr := postgres.NewGroupRepo(sqlDB)

	groups, err := gr.GetAllGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(groups) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No groups found"})
		return
	}

	var groupsResponse []response_models.GetGroupResponse
	for _, groupModel := range groups {
		groupsResponse = append(groupsResponse, response_models.FormatGetGroupResponse(groupModel))
	}

	c.JSON(http.StatusOK, groupsResponse)
}
