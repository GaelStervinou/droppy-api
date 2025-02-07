package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/http/response_models"
	"go-api/internal/repositories"
	commentservice "go-api/internal/services/comment"
	pushnotificationservice "go-api/internal/services/push_notification"
	"go-api/pkg/model"
	"log"
	"net/http"
	"strconv"
)

// CommentDrop godoc
//
//	@Summary		Comment on a drop
//	@Description	Comment on a drop
//	@Tags			drop
//	@Accept			json
//	@Produce		json
//
// @Security BearerAuth
//
//	@Param			id path int true "Drop ID"
//	@Param			comment	body		model.CommentCreationParam	true	"Comment creation object"
//	@Success		201	{object} response_models.GetCommentResponse
//	@Failure		401
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Router			/drops/{id}/comments [post]
func CommentDrop(c *gin.Context) {
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

	dropIdUint, err := strconv.ParseUint(dropId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid drop ID"})
		return
	}

	var commentCreationParam model.CommentCreationParam

	if err := c.ShouldBindJSON(&commentCreationParam); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	cs := &commentservice.CommentService{
		Repo: repositories.Setup(),
	}

	comment, err := cs.CommentDrop(uint(dropIdUint), uintCurrentUserId, commentCreationParam)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := response_models.FormatGetCommentResponse(comment)

	c.JSON(http.StatusCreated, response)

	drop, err := cs.Repo.DropRepository.GetDropById(uint(dropIdUint))

	if err != nil {
		return
	}

	user, err := cs.Repo.UserRepository.GetById(drop.GetCreatedById())

	if err != nil {
		return
	}

	if user.GetFCMToken() != "" {
		pushNotificationService := &pushnotificationservice.PushNotificationService{Repo: repositories.Setup()}
		err = pushNotificationService.SendNotification("like", []string{user.GetFCMToken()})
		if err != nil {
			log.Printf("Error: Error sending push notification: %v", err)
		}
	}

	followers, err := cs.Repo.FollowRepository.GetFollowers(drop.GetCreatedById())

	if err != nil {
		return
	}

	for _, follower := range followers {
		_ = NewDropAvailable(follower.GetFollowerID(), drop)
	}

	_ = NewDropAvailable(drop.GetCreatedById(), drop)
}

// DeleteComment godoc
//
//	@Summary		Delete a comment
//	@Description	Delete a comment
//	@Tags			drop
//	@Accept			json
//	@Produce		json
//
// @Security BearerAuth
//
//	@Param			id path int true "Comment ID"
//	@Success		204 No Content
//	@Failure		401
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Failure		500
//	@Router			/comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	currentUserId, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	commentId := c.Param("id")

	if "" == commentId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	commentIdUint, err := strconv.ParseUint(commentId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	cs := &commentservice.CommentService{
		Repo: repositories.Setup(),
	}

	err = cs.CanDeleteComment(uint(commentIdUint), currentUserId.(uint))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	comment, err := cs.Repo.CommentRepository.GetById(uint(commentIdUint))

	if err != nil {
		return
	}

	err = cs.DeleteComment(uint(commentIdUint))

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)

	drop, err := cs.Repo.DropRepository.GetDropById(comment.GetDrop().GetID())

	if err != nil {
		return
	}

	followers, err := cs.Repo.FollowRepository.GetFollowers(drop.GetCreatedById())

	if err != nil {
		return
	}

	for _, follower := range followers {
		_ = NewDropAvailable(follower.GetFollowerID(), drop)
	}

	_ = NewDropAvailable(drop.GetCreatedById(), drop)
}
