package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/repositories"
	dropservice "go-api/internal/services/drop"
	"go-api/internal/storage/postgres"
	"go-api/internal/storage/postgres/drop"
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
//	@Success		201	{object} drop.Drop
//	@Failure		422
//	@Failure		500
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

	c.JSON(http.StatusCreated, createdDrop)
}

// DropsByUserId godoc
//
//	@Summary		List drops by user ID
//	@Description	List drops by user ID
//	@Tags			drop
//	@Accept			json
//	@Produce		json
//	@Param			id path int true "User ID"
//	@Success		200	{object} []drop.Drop
//	@Failure		422
//	@Failure		500
//	@Router			/users/:id/drops [get]
func DropsByUserId(c *gin.Context) {
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

	dr := drop.NewRepo(sqlDB)

	drops, err := dr.GetUserDrops(uintCurrentUserId)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, drops)
}

// GetCurrentUserFeed godoc
//
//	@Summary		Get feed
//	@Description	Get feed
//	@Tags			user
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
