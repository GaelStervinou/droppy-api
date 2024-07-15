package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-api/internal/http/response_models"
	"go-api/internal/repositories"
	"go-api/internal/services/follow"
	"go-api/internal/storage/postgres"
	"go-api/pkg/converters"
	"go-api/pkg/errors2"
	"go-api/pkg/model"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// FollowUser godoc
//
// @Summary		Follow user
// @Description	Follow user by its ID
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security BearerAuth
//
//	@Param			user	body		model.FollowCreationParam	true	"Follow creation object"
//
// @Success		201	{object} postgres.Follow
// @Failure		422
// @Failure		401
// @Failure		500
// @Router			/follows [post]
func FollowUser(c *gin.Context) {
	var followCreationParam model.FollowCreationParam

	if err := c.ShouldBindJSON(&followCreationParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	sqlDB := postgres.Connect()

	us := postgres.NewUserRepo(sqlDB)

	requestedUser, err := us.GetById(followCreationParam.UserToFollowID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if nil == requestedUser {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "User not found"})
		return
	}

	if uintCurrentUserId == followCreationParam.UserToFollowID {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "You can't follow yourself"})
		return
	}

	isFollowingAllowed, err := us.CanUserBeFollowed(followCreationParam.UserToFollowID)

	if err != nil || !isFollowingAllowed {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "You can't follow this user for now. Please try again later"})
		return
	}

	followRepo := postgres.NewFollowRepo(sqlDB)

	if alreadyFollowing, _ := followRepo.AreAlreadyFollowing(uintCurrentUserId, followCreationParam.UserToFollowID); alreadyFollowing {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "You are already following this user"})
		return
	}

	createdFollow, err := followRepo.Create(uintCurrentUserId, followCreationParam.UserToFollowID, !requestedUser.IsPrivateUser())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if nil == createdFollow {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not followed"})
		return
	}

	if createdFollow.GetStatus() == new(postgres.FollowPendingStatus).ToInt() {
		fmt.Printf("Sending pending follows for user %d\n", uintCurrentUserId)
		err = SendPendingFollowsWS(followCreationParam.UserToFollowID, postgres.NewFollowRepo(sqlDB))
		if err != nil {
			log.Printf("Error sending message to user %d: %v", uintCurrentUserId, err)
			return
		}
	} else if createdFollow.GetStatus() == new(postgres.FollowAcceptedStatus).ToInt() {
		fmt.Printf("Sending drops to user %d\n", uintCurrentUserId)
		NewDropAvailable(uintCurrentUserId)
	}

	c.JSON(http.StatusCreated, createdFollow)
}

var pendingFollowUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type PendingFollowWebSocketConnection struct {
	conn *websocket.Conn
}

var userPendingFollowConnections = make(map[string]*PendingFollowWebSocketConnection)
var muPendingFollow sync.Mutex

func GetMyPendingRequestsWS(c *gin.Context) {
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

	conn, err := pendingFollowUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade WebSocket"})
		return
	}

	wsConn := &PendingFollowWebSocketConnection{conn: conn}

	muPendingFollow.Lock()
	userPendingFollowConnections[strconv.Itoa(int(uintCurrentUserId))] = wsConn
	fmt.Printf("Users connected: %v\n", userPendingFollowConnections)
	muPendingFollow.Unlock()

	sqlDB := postgres.Connect()

	err = SendPendingFollowsWS(uintCurrentUserId, postgres.NewFollowRepo(sqlDB))
	if err != nil {
		log.Printf("Error sending message to user %d: %v", uintCurrentUserId, err)
		return
	}

	defer func() {
		muPendingFollow.Lock()
		delete(userPendingFollowConnections, strconv.Itoa(int(uintCurrentUserId)))
		muPendingFollow.Unlock()
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing WebSocket connection: %v", err)
		}
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// AcceptRequest godoc
//
// @Summary		 Accept request
// @Description	 Accept request by its ID
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security BearerAuth
//
// @Success		201	{object} []postgres.Follow
// @Failure		422
// @Failure		401
// @Failure		500
// @Router			/follows/accept/{id} [post]
func AcceptRequest(c *gin.Context) {
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

	sqlDB := postgres.Connect()

	id := c.Param("id")

	if "" == id {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	followId, err := converters.StringToUint(id)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid follow ID"})
		return
	}

	followRepo := postgres.NewFollowRepo(sqlDB)

	myFollow, err := followRepo.GetFollowByID(followId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if myFollow.GetFollowed().GetID() != uintCurrentUserId || myFollow.GetStatus() != new(postgres.FollowPendingStatus).ToInt() {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	err = followRepo.AcceptRequest(followId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = SendPendingFollowsWS(uintCurrentUserId, followRepo)
	if err != nil {
		log.Printf("Error sending message to user %d: %v", uintCurrentUserId, err)
	}

	acceptedFollow, err := followRepo.GetFollowByID(followId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	NewDropAvailable(acceptedFollow.GetFollowerID())

	c.JSON(http.StatusOK, gin.H{"message": "Follow request accepted"})
}

// RejectRequest godoc
//
// @Summary		 Refuse follow request
// @Description	 Refuse follow request by its ID
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security BearerAuth
//
// @Success		201
// @Failure		422
// @Failure		401
// @Failure		500
// @Router			/follows/refuse/{id} [delete]
func RejectRequest(c *gin.Context) {
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

	sqlDB := postgres.Connect()

	id := c.Param("id")

	if "" == id {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	followId, err := converters.StringToUint(id)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid follow ID"})
		return
	}

	followRepo := postgres.NewFollowRepo(sqlDB)

	myFollow, err := followRepo.GetFollowByID(followId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if myFollow.GetFollowed().GetID() != uintCurrentUserId || myFollow.GetStatus() != new(postgres.FollowPendingStatus).ToInt() {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	err = followRepo.RejectRequest(followId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = SendPendingFollowsWS(uintCurrentUserId, followRepo)
	if err != nil {
		log.Printf("Error sending message to user %d: %v", uintCurrentUserId, err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Follow request refused"})
}

// GetUserFollowing godoc
//
// @Summary		Get user following
// @Description	Get user following
// @Tags			user
// @Accept			json
// @Produce		json
// @Security BearerAuth
//
// @Success		200	{object} []response_models.GetUserResponseInterface
// @Failure		422
// @Failure		403
// @Failure		401
// @Failure		500
// @Router			/users/{id}/following [get]
func GetUserFollowing(c *gin.Context) {
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

	userID := c.Param("id")
	if "" == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	userIDuint, err := converters.StringToUint(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fs := &follow.FollowService{
		Repo: repositories.Setup(),
	}

	following, err := fs.GetUserFollowing(userIDuint, uintCurrentUserId)

	if err != nil {
		var notAllowedErr errors2.NotAllowedError
		if errors.As(err, &notAllowedErr) {
			c.JSON(http.StatusForbidden, gin.H{"error": notAllowedErr.Reason})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followingResponse []response_models.GetUserResponseInterface
	for _, followResp := range following {
		followingResponse = append(followingResponse, response_models.FormatGetUserResponse(followResp.GetFollowed()))
	}

	c.JSON(http.StatusOK, followingResponse)
}

// GetUserFollowers godoc
//
// @Summary		Get user followers
// @Description	Get user followers
// @Tags			user
// @Accept			json
// @Produce		json
// @Security BearerAuth
//
// @Success		200	{object} []response_models.GetUserResponseInterface
// @Failure		422
// @Failure		403
// @Failure		401
// @Failure		500
// @Router			/users/{id}/followers [get]
func GetUserFollowers(c *gin.Context) {
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

	userID := c.Param("id")
	if "" == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	userIDuint, err := converters.StringToUint(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fs := &follow.FollowService{
		Repo: repositories.Setup(),
	}

	followers, err := fs.GetUserFollowers(userIDuint, uintCurrentUserId)

	if err != nil {
		var notAllowedErr errors2.NotAllowedError
		if errors.As(err, &notAllowedErr) {
			c.JSON(http.StatusForbidden, gin.H{"error": notAllowedErr.Reason})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followersResponse []response_models.GetUserResponseInterface
	for _, followResp := range followers {
		followersResponse = append(followersResponse, response_models.FormatGetUserResponse(followResp.GetFollower()))
	}

	c.JSON(http.StatusOK, followersResponse)
}

func SendPendingFollowsWS(userID uint, followRepo model.FollowRepository) error {
	fmt.Printf("Sending pending follows for user %d\n", userID)
	muPendingFollow.Lock()
	wsConn, ok := userPendingFollowConnections[strconv.Itoa(int(userID))]
	muPendingFollow.Unlock()
	if !ok {
		fmt.Printf("Connections are: %v\n", userPendingFollowConnections)
		fmt.Printf("User %d not found in userPendingFollowConnections\n", userID)
		return nil
	}

	pendingRequests, err := followRepo.GetPendingRequests(userID)
	if err != nil {
		return err
	}

	var pendingFollowResponses []response_models.GetOnePendingFollowResponse
	for _, pendingFollow := range pendingRequests {
		pendingFollowResponses = append(pendingFollowResponses, response_models.FormatGetOnePendingFollowResponse(pendingFollow))
	}

	return wsConn.conn.WriteJSON(pendingFollowResponses)
}

// DeleteFollow godoc
//
// @Summary		Delete follow
// @Description	Delete follow
// @Tags			follow
// @Accept			json
// @Produce		json
// @Param			id path int true "Follow ID"
// @Security BearerAuth
//
// @Success		204 No Content
// @Failure		422
// @Failure		401
// @Failure		403
// @Failure		500
// @Router			/follows/{id} [delete]
func DeleteFollow(c *gin.Context) {
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
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	followId, err := converters.StringToUint(id)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid follow ID"})
		return
	}

	fs := &follow.FollowService{
		Repo: repositories.Setup(),
	}

	err = fs.DeleteFollow(uintCurrentUserId, followId)

	if err != nil {
		var notAllowedErr errors2.NotAllowedError
		if errors.As(err, &notAllowedErr) {
			c.JSON(http.StatusForbidden, gin.H{"error": notAllowedErr.Reason})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
