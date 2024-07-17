package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-api/internal/http/response_models"
	"go-api/internal/repositories"
	pushnotificationservice "go-api/internal/services/push_notification"
	reportservice "go-api/internal/services/report"
	"go-api/internal/storage/postgres"
	"go-api/pkg/model"
	"net/http"
	"strconv"
	"strings"
)

// GetAllUsers godoc
//
// @Summary		Get all users
// @Description	Get all users by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			page query int false "Page number"
// @Param			pageSize query int false "Page size"
// @Success		200	{object} []response_models.GetUserResponse
// @Failure		500
// @Router			/admin/users [get]
func GetAllUsers(c *gin.Context) {
	sqlDB := postgres.Connect()

	us := postgres.NewUserRepo(sqlDB)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	users, err := us.GetAll(page, pageSize)
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

// GetAllUsersCount godoc
//
// @Summary		Get all users count
// @Description	Get all users count by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Success		200
// @Failure		500
// @Router			/admin/users/count [get]
func GetAllUsersCount(c *gin.Context) {
	sqlDB := postgres.Connect()

	us := postgres.NewUserRepo(sqlDB)

	count, err := us.GetAllUserCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, count)
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
// @Param		user	body model.AdminUpdateUserRequest true "User data"
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

	users, err := us.GetAll(1, 20)
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

// GetAllGroups godoc
//
// @Summary		Get all groups
// @Description	Get all groups by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			page query int false "Page number"
// @Param			pageSize query int false "Page size"
// @Success		200	{object} []response_models.GetGroupResponse
// @Failure		500
// @Router			/admin/groups [get]
func GetAllGroups(c *gin.Context) {
	sqlDB := postgres.Connect()

	gr := postgres.NewGroupRepo(sqlDB)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	groups, err := gr.GetAllGroups(page, pageSize)
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

// GetAllGroupsCount godoc
//
// @Summary		Get all groups count
// @Description	Get all groups count by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Success		200
// @Failure		500
// @Router			/admin/groups/count [get]
func GetAllGroupsCount(c *gin.Context) {
	sqlDB := postgres.Connect()

	gr := postgres.NewGroupRepo(sqlDB)

	count, err := gr.GetAllGroupsCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, count)

}

// GetAllDrops godoc
//
// @Summary		Get all drops
// @Description	Get all drops by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			page query int false "Page number"
// @Param			pageSize query int false "Page size"
// @Success		200	{object} []response_models.GetDropResponse
// @Failure		500
// @Router			/admin/drops [get]
func GetAllDrops(c *gin.Context) {
	sqlDB := postgres.Connect()

	dr := postgres.NewDropRepo(sqlDB)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	drops, err := dr.GetAllDrops(page, pageSize)
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

// GetAllDropsCount godoc
//
// @Summary		Get all drops count
// @Description	Get all drops count by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Success		200
// @Failure		500
// @Router			/admin/drops/count [get]
func GetAllDropsCount(c *gin.Context) {
	sqlDB := postgres.Connect()

	dr := postgres.NewDropRepo(sqlDB)

	count, err := dr.GetAllDropsCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, count)
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
// @Param			page query int false "Page number"
// @Param			pageSize query int false "Page size"
// @Success		200
// @Failure		500
// @Router			/admin/reports [get]
func GetAllReports(c *gin.Context) {
	sqlDB := postgres.Connect()

	rr := postgres.NewReportRepo(sqlDB)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	fmt.Printf("page: %d, pageSize: %d\n", page, pageSize)

	reports, err := rr.GetAllReports(page, pageSize)
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

// AdminDeleteDrop godoc
//
// @Summary		Delete drop
// @Description	Delete drop by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			id path string true "Drop ID"
// @Success		200
// @Failure		500
// @Router			/admin/drops/{id} [delete]
func AdminDeleteDrop(c *gin.Context) {
	sqlDB := postgres.Connect()

	dr := postgres.NewDropRepo(sqlDB)

	dropID := c.Param("id")
	uintDropID, err := strconv.ParseUint(dropID, 10, 64)

	err = dr.Delete(uint(uintDropID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	drops, err := dr.GetAllDrops(1, 20)
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

// AdminDeleteGroup godoc
//
// @Summary		Delete group
// @Description	Delete group by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			id path string true "Group ID"
// @Success		200
// @Failure		500
// @Router			/admin/groups/{id} [delete]
func AdminDeleteGroup(c *gin.Context) {
	sqlDB := postgres.Connect()

	gr := postgres.NewGroupRepo(sqlDB)

	groupID := c.Param("id")
	uintGroupID, err := strconv.ParseUint(groupID, 10, 64)

	err = gr.Delete(uint(uintGroupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groups, err := gr.GetAllGroups(1, 20)
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

// AdminDeleteComment godoc
//
// @Summary		Delete comment
// @Description	Delete comment by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			id path string true "Comment ID"
// @Success		200
// @Failure		500
// @Router			/admin/comments/{id} [delete]
func AdminDeleteComment(c *gin.Context) {
	sqlDB := postgres.Connect()

	cr := postgres.NewCommentRepo(sqlDB)

	commentID := c.Param("id")
	uintCommentID, err := strconv.ParseUint(commentID, 10, 64)

	err = cr.DeleteComment(uint(uintCommentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

// AdminManageReport godoc
//
// @Summary		Manage report
// @Description	Manage report by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			id path string true "Report ID"
// @Param 		report body model.ManageReportRequest true "Manage report data"
// @Success		200
// @Failure		500
// @Router			/admin/reports/{id} [put]
func AdminManageReport(c *gin.Context) {
	sqlDB := postgres.Connect()

	rr := postgres.NewReportRepo(sqlDB)

	reportID := c.Param("id")

	var manageReportRequest model.ManageReportRequest
	if err := c.ShouldBindJSON(&manageReportRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uintReportID, err := strconv.ParseUint(reportID, 10, 64)
	reportModel, err := rr.GetReportById(uint(uintReportID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reportService := reportservice.ReportService{
		Repo: repositories.Setup(),
	}

	reportModel, err = reportService.ManageReport(reportModel.GetID(), manageReportRequest.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reports, err := rr.GetAllReports(1, 20)
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

// AdminScheduleDrop godoc
//
// @Summary		Schedule drop
// @Description	Schedule drop by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param		params	body model.ScheduleDropParam true "Send drop data"
// @Success		201 {object} response_models.GetDropNotificationResponse
// @Failure		422 "Invalid drop data"
// @Failure		500
// @Router			/admin/drops/schedule [post]
func AdminScheduleDrop(c *gin.Context) {
	sqlDB := postgres.Connect()

	dnr := postgres.NewDropNotifRepo(sqlDB)

	var scheduleDropParam model.ScheduleDropParam
	if err := c.ShouldBindJSON(&scheduleDropParam); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	dropNotifModel, err := dnr.Create(scheduleDropParam.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response_models.FormatGetDropNotificationResponse(dropNotifModel))
}

// AdminSendDropNow godoc
//
// @Summary		Send drop now
// @Description	Send drop now by admin user
// @Tags			admin
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			params	body model.SendDropRequest true "Send drop data"
// @Success		201 {object} response_models.GetDropNotificationResponse
// @Failure		422 "Invalid drop data"
// @Failure		500
// @Router			/admin/drops/send-now [post]
func AdminSendDropNow(c *gin.Context) {
	sqlDB := postgres.Connect()

	dnr := postgres.NewDropNotifRepo(sqlDB)

	var scheduleDropParam model.ScheduleDropParam
	if err := c.ShouldBindJSON(&scheduleDropParam); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	dropType := strings.ToLower(scheduleDropParam.Type)
	dropNotifModel, err := dnr.Create(dropType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pushNotificationService := pushnotificationservice.PushNotificationService{
		Repo: repositories.Setup(),
	}

	pushNotificationService.SendNotificationsToAllUser(scheduleDropParam.Type)

	RefreshHasUserDroppedToday()

	c.JSON(http.StatusCreated, response_models.FormatGetDropNotificationResponse(dropNotifModel))
}
