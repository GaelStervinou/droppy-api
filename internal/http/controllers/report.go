package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/repositories"
	"go-api/internal/services/report"
	"go-api/pkg/model"
	"net/http"
)

// CreateReport godoc
//
// @Summary		Create report
// @Description	Create report
// @Tags			report
// @Accept			json
// @Produce		json
// @Security BearerAuth
// @Param			report	body		model.ReportCreationParam	true	"Report object"
// @Success		201	{object} response_models.GetReportResponse
// @Failure		422 {object} errors2.MultiFieldsError
// @Router			/reports [post]
func CreateReport(c *gin.Context) {
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

	var reportRequest model.ReportCreationParam
	if err := c.ShouldBindJSON(&reportRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reportService := report.ReportService{
		Repo: repositories.Setup(),
	}

	createdReport, err := reportService.CreateReport(uintCurrentUserId, reportRequest)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdReport)
}
