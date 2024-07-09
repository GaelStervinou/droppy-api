package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/http/response_models"
	"go-api/internal/repositories"
	commentresponseservice "go-api/internal/services/comment_response"
	"go-api/pkg/model"
	"net/http"
	"strconv"
)

// RespondToComment godoc
//
//	@Summary		Respond to a comment
//	@Description	Respond to a comment
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//
// @Security BearerAuth
//
//	@Param			id path int true "Comment ID"
//	@Param			content body string true "Comment content"
//	@Success		201	{object} response_models.GetCommentResponseResponse
//	@Failure		401
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Router			/comments/{id}/responses [post]
func RespondToComment(c *gin.Context) {
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

	var commentResponseCreationParam model.CommentCreationParam

	if err := c.ShouldBindJSON(&commentResponseCreationParam); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	cs := &commentresponseservice.CommentResponseService{
		Repo: repositories.Setup(),
	}

	commentResponse, err := cs.RespondToComment(uint(commentIdUint), uintCurrentUserId, commentResponseCreationParam)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := response_models.FormatGetCommentResponseResponse(commentResponse)

	c.JSON(http.StatusCreated, response)
}

// DeleteCommentResponse godoc
//
//	@Summary		Delete a comment response
//	@Description	Delete a comment response
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//
// @Security BearerAuth
//
//	@Param			id path int true "Comment ID"
//	@Param			responseId path int true "Comment response ID"
//	@Success		204	{} No Content
//	@Failure		401
//	@Failure		422 {object} errors2.MultiFieldsError
//	@Router			/comments/{id}/responses/{responseId} [delete]
func DeleteCommentResponse(c *gin.Context) {
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

	commentResponseId := c.Param("responseId")

	if "" == commentResponseId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	commentResponseIdUint, err := strconv.ParseUint(commentResponseId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment response ID"})
		return
	}

	cs := &commentresponseservice.CommentResponseService{
		Repo: repositories.Setup(),
	}

	err = cs.DeleteCommentResponse(uint(commentResponseIdUint), uintCurrentUserId)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
