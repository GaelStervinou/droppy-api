package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/http/response_models"
	"go-api/internal/storage/postgres"
	"go-api/pkg/model"
	"net/http"
	"strconv"
)

// GetAllUsers godoc
//
// @Summary		Get all users
// @Description	Get all users by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Success		200	{object} []response_models.GetUserResponse
// @Failure		500
// @Router			/admin/users [get]
func GetAllUsers(c *gin.Context) {
	sqlDB := postgres.Connect()

	us := postgres.NewUserRepo(sqlDB)

	users, err := us.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var usersResponse []response_models.GetUserResponseInterface
	for _, userModel := range users {
		usersResponse = append(usersResponse, response_models.FormatAdminGetUserResponse(userModel))
	}

	c.JSON(http.StatusOK, usersResponse)
}

// UpdateUser godoc
//
// @Summary		Update user
// @Description	Update user by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			id path string true "User ID"
// @Param			body request_models.UpdateUserRequest true "User data"
// @Success		200
// @Failure		400
// @Failure		500
// @Router			/admin/users/{id} [put]
func UpdateUser(c *gin.Context) {
	sqlDB := postgres.Connect()

	us := postgres.NewUserRepo(sqlDB)

	userID := c.Param("id")

	var updateUserRequest model.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&updateUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uintUserID, err := strconv.ParseUint(userID, 10, 64)
	userModel, err := us.GetById(uint(uintUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userModel, err = us.UpdateByAdmin(userModel.GetID(), updateUserRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response_models.FormatAdminGetUserResponse(userModel))
}

// BanUser godoc
//
// @Summary		Ban user
// @Description	Ban user by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			id path string true "User ID"
// @Success		200
// @Failure		500
// @Router			/admin/users/{id}/ban [put]
func BanUser(c *gin.Context) {
	sqlDB := postgres.Connect()

	us := postgres.NewUserRepo(sqlDB)

	userID := c.Param("id")
	uintUserID, err := strconv.ParseUint(userID, 10, 64)

	userModel, err := us.GetById(uint(uintUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userModel, err = us.BanUser(userModel.GetID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response_models.FormatAdminGetUserResponse(userModel))
}

// UnbanUser godoc
//
// @Summary		Unban user
// @Description	Unban user by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			id path string true "User ID"
// @Success		200
// @Failure		500
// @Router			/admin/users/{id}/unban [put]
func UnbanUser(c *gin.Context) {
	sqlDB := postgres.Connect()

	us := postgres.NewUserRepo(sqlDB)

	userID := c.Param("id")
	uintUserID, err := strconv.ParseUint(userID, 10, 64)

	userModel, err := us.GetById(uint(uintUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userModel, err = us.UnbanUser(userModel.GetID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response_models.FormatAdminGetUserResponse(userModel))
}

// GetAllGroups godoc
//
// @Summary		Get all groups
// @Description	Get all groups by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Success		200	{object} []response_models.GetGroupResponse
// @Failure		500
// @Router			/admin/groups [get]
func GetAllGroups(c *gin.Context) {
	sqlDB := postgres.Connect()

	gr := postgres.NewGroupRepo(sqlDB)

	groups, err := gr.GetAllGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var groupsResponse []response_models.GetGroupResponse
	for _, groupModel := range groups {
		groupsResponse = append(groupsResponse, response_models.FormatGetGroupResponse(groupModel))
	}

	c.JSON(http.StatusOK, groupsResponse)
}

// GetAllDrops godoc
//
// @Summary		Get all drops
// @Description	Get all drops by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Success		200	{object} []response_models.GetDropResponse
// @Failure		500
// @Router			/admin/drops [get]
func GetAllDrops(c *gin.Context) {
	sqlDB := postgres.Connect()

	dr := postgres.NewDropRepo(sqlDB)

	drops, err := dr.GetAllDrops()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dropsResponse []response_models.GetDropResponse
	for _, dropModel := range drops {
		dropsResponse = append(dropsResponse, response_models.FormatGetDropResponse(dropModel, false))
	}

	c.JSON(http.StatusOK, dropsResponse)
}

// GetAllComments godoc
//
// @Summary		Get all comments
// @Description	Get all comments by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Success		200	{object} []response_models.GetCommentResponse
// @Failure		500
// @Router			/admin/comments [get]
func GetAllComments(c *gin.Context) {
	sqlDB := postgres.Connect()

	cr := postgres.NewCommentRepo(sqlDB)

	comments, err := cr.GetAllComments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var commentsResponse []response_models.GetCommentResponse
	for _, commentModel := range comments {
		commentsResponse = append(commentsResponse, response_models.FormatGetCommentResponse(commentModel))
	}

	c.JSON(http.StatusOK, commentsResponse)
}

// GetAllReports godoc
//
// @Summary		Get all reports
// @Description	Get all reports by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Success		200
// @Failure		500
// @Router			/admin/reports [get]
func GetAllReports(c *gin.Context) {
	sqlDB := postgres.Connect()

	rr := postgres.NewReportRepo(sqlDB)

	reports, err := rr.GetAllReports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var reportsResponse []response_models.GetReportResponse
	for _, reportModel := range reports {
		reportsResponse = append(reportsResponse, response_models.FormatGetReportResponse(reportModel))
	}

	c.JSON(http.StatusOK, reportsResponse)
}
