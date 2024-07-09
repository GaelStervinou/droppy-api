package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/repositories"
	likeservice "go-api/internal/services/like"
	"go-api/pkg/model"
	"net/http"
	"strconv"
)

// LikeDrop godoc
//
//	@Summary		Like Drop
//	@Description	Like Drop
//	@Tags			drop
//	@Accept			json
//	@Produce		json
//
// @Security BearerAuth
//
//	@Param			id path int true "Drop ID"
//	@Success		200	{object} response_models.GetDropResponse
//	@Failure		401
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		500
//	@Router			/drops/{id}/like [post]
func LikeDrop(c *gin.Context) {
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

	dropId := c.Param("id")

	if "" == dropId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	uintDropId, err := strconv.ParseUint(dropId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid drop ID"})
		return
	}

	likeParam := model.LikeParam{DropId: uint(uintDropId)}

	ls := &likeservice.LikeService{Repo: repositories.Setup()}

	like, err := ls.LikeDrop(uintCurrentUserId, likeParam)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, like)
}

// UnlikeDrop godoc
//
//	@Summary		Unlike Drop
//	@Description	Unlike Drop
//	@Tags			drop
//	@Accept			json
//	@Produce		json
//
// @Security BearerAuth
//
//	@Param			id path int true "Drop ID"
//	@Success		204	No Content
//	@Failure		401
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		500
//	@Router			/drops/{id}/unlike [delete]
func UnlikeDrop(c *gin.Context) {
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

	dropId := c.Param("id")

	if "" == dropId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	uintDropId, err := strconv.ParseUint(dropId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid drop ID"})
		return
	}

	likeParam := model.LikeParam{DropId: uint(uintDropId)}

	ls := &likeservice.LikeService{Repo: repositories.Setup()}

	err = ls.UnlikeDrop(uintCurrentUserId, likeParam)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
