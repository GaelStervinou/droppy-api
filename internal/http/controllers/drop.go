package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/http/response_models"
	"go-api/internal/repositories"
	dropservice "go-api/internal/services/drop"
	"go-api/internal/storage/postgres"
	"go-api/internal/storage/postgres/user"
	"go-api/pkg/converters"
	"go-api/pkg/model"
	"net/http"
)

// CreateDrop godoc
//
//	@Summary		Create a new drop
//	@Description	Create a new drop
//	@Tags			drop
//	@Accept			json
//
// @Security BearerAuth
//
//	@Produce		json
//	@Param			drop body		model.DropCreationParam	true	"Drop object"
//	@Success		201	{object} response_models.GetDropResponse
//	@Failure		401
//	@Failure		422
//	@Router			/drops [post]
func CreateDrop(c *gin.Context) {
	var dropCreationParam model.DropCreationParam

	if err := c.ShouldBindJSON(&dropCreationParam); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
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

	ds := &dropservice.DropService{
		Repo: repositories.Setup(),
	}

	createdDrop, err := ds.CreateDrop(uintCurrentUserId, dropCreationParam)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	sqlDB, err := postgres.Connect()
	us := user.NewRepo(sqlDB)

	createdBy, err := us.GetById(uintCurrentUserId)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	createdByResponse := response_models.FormatGetUserResponse(createdBy)

	response := response_models.FormatGetDropResponse(createdDrop, createdByResponse)

	c.JSON(http.StatusCreated, response)
}

// DropsByUserId godoc
//
//	@Summary		List drops by user ID
//	@Description	List drops by user ID
//	@Tags			drop
//	@Accept			json
//	@Produce		json
//	@Param			id path int true "User ID"
//	@Success		200	{object} []response_models.GetDropResponse
//	@Failure		500
//	@Router			/users/:id/drops [get]
func DropsByUserId(c *gin.Context) {
	repo := repositories.Setup()

	ds := &dropservice.DropService{
		Repo: repo,
	}

	currentUserId, exists := c.Get("userId")
	var currentUser model.UserModel

	if exists {
		uintCurrentUserId, ok := currentUserId.(uint)
		if ok {

			targetedUser, err := repo.UserRepository.GetById(uintCurrentUserId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			currentUser = targetedUser
		}
	}

	id := c.Param("id")

	if "" == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	userId, err := converters.StringToUint(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	drops, err := ds.GetDropsByUserId(userId, currentUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dropsResponse []response_models.GetDropResponse
	dropUser, err := repo.UserRepository.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponse := response_models.FormatGetUserResponse(dropUser)
	for _, drop := range drops {
		dropResponse := response_models.FormatGetDropResponse(drop, userResponse)
		dropsResponse = append(dropsResponse, dropResponse)
	}

	c.JSON(200, dropsResponse)
}

// GetCurrentUserFeed godoc
//
//	@Summary		Get feed
//	@Description	Get feed
//	@Tags			drop
//	@Accept			json
//	@Produce		json
//	@Security BearerAuth
//	@Success		200	{object} []drop.Drop
//	@Failure		500
//	@Router			/users/my-feed [get]
func GetCurrentUserFeed(c *gin.Context) {
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

	ds := &dropservice.DropService{
		Repo: repositories.Setup(),
	}

	drops, err := ds.GetUserFeed(uintCurrentUserId)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, drops)
}

type GetUserDropParam struct {
	Id uint `uri:"id" binding:"required"`
}
