package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-api/internal/http/response_models"
	"go-api/internal/repositories"
	dropservice "go-api/internal/services/drop"
	"go-api/internal/storage/postgres"
	"go-api/pkg/converters"
	"go-api/pkg/drop_type_apis"
	"go-api/pkg/model"
	"log"
	"net/http"
	"strconv"
	"sync"
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

	response := response_models.FormatGetDropResponse(createdDrop, false)

	c.JSON(http.StatusCreated, response)

	sqlDB := postgres.Connect()

	fr := postgres.NewFollowRepo(sqlDB)

	userFollowers, err := fr.GetFollowers(uintCurrentUserId)

	if err != nil {
		return
	}

	fmt.Printf("Sending new drop to followers %v\n", userFollowers)
	for _, follower := range userFollowers {
		NewDropAvailable(follower.GetFollowerID())
	}
}

func GetOneDrop(c *gin.Context) {
	repo := repositories.Setup()

	ds := &dropservice.DropService{
		Repo: repo,
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

	id := c.Param("id")

	if "" == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	dropId, err := converters.StringToUint(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	drop, err := ds.GetDropById(dropId, uintCurrentUserId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	isCurrentUserLiking, err := ds.IsCurrentUserLiking(drop.GetID(), uintCurrentUserId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dropResponse := response_models.FormatGetDropResponse(drop, isCurrentUserLiking)

	c.JSON(200, dropResponse)
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

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uintCurrentUserId, ok := currentUserId.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser, err := repo.UserRepository.GetById(uintCurrentUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	for _, drop := range drops {
		isCurrentUserLiking, err := ds.IsCurrentUserLiking(drop.GetID(), uintCurrentUserId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		dropResponse := response_models.FormatGetDropResponse(drop, isCurrentUserLiking)
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
//	@Success		200	{object} []postgres.Drop
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

// Define the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Create a struct to manage WebSocket connections
type WebSocketConnection struct {
	conn *websocket.Conn
}

// Global map to store connections per user
var userFeedConnections = make(map[string]*WebSocketConnection)
var mu sync.Mutex

func GetCurrentUserFeedWS(c *gin.Context) {
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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade WebSocket"})
		return
	}

	wsConn := &WebSocketConnection{conn: conn}

	mu.Lock()
	userFeedConnections[strconv.Itoa(int(uintCurrentUserId))] = wsConn
	fmt.Printf("Users connected: %v\n", userFeedConnections)
	mu.Unlock()

	ds := &dropservice.DropService{
		Repo: repositories.Setup(),
	}

	/*hasDropped, err := ds.HasUserDroppedToday(uintCurrentUserId)

	if err != nil {
		log.Printf("Error checking if user has dropped today: %v", err)
		return
	}*/

	availableDrops, err := ds.GetUserFeed(uintCurrentUserId)

	if err != nil {
		log.Printf("Error getting user feed: %v", err)
		return
	}

	var dropResponses []response_models.GetDropResponse
	for _, drop := range availableDrops {
		/*if !hasDropped {
			//TODO ne pas envoyer la pic, le content et la description ( donc fair eun interface pour les 2 types de drop et déclarer une var au dessus de ce type là )
		}*/
		isCurrentUserLiking, err := ds.IsCurrentUserLiking(drop.GetID(), uintCurrentUserId)

		if err != nil {
			log.Printf("Error checking if user is liking drop: %v", err)
			return
		}

		dropResponse := response_models.FormatGetDropResponse(drop, isCurrentUserLiking)
		dropResponses = append(dropResponses, dropResponse)
	}

	err = wsConn.conn.WriteJSON(dropResponses)
	if err != nil {
		log.Printf("Error sending message to user %d: %v", uintCurrentUserId, err)
		return
	}

	defer func() {
		mu.Lock()
		delete(userFeedConnections, strconv.Itoa(int(uintCurrentUserId)))
		mu.Unlock()
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

func sendToUser(userID uint, ds *dropservice.DropService) error {
	mu.Lock()
	wsConn, ok := userFeedConnections[strconv.Itoa(int(userID))]
	fmt.Printf("Users: %v\n", userFeedConnections)
	mu.Unlock()

	if !ok {
		return fmt.Errorf("user not connected")
	}

	availableDrops, err := ds.GetUserFeed(userID)

	if err != nil {
		log.Printf("Error getting user feed: %v", err)
		return err
	}

	var dropResponses []response_models.GetDropResponse
	for _, drop := range availableDrops {
		isCurrentUserLiking, err := ds.IsCurrentUserLiking(drop.GetID(), userID)

		if err != nil {
			log.Printf("Error checking if user is liking drop: %v", err)
			return err
		}

		dropResponse := response_models.FormatGetDropResponse(drop, isCurrentUserLiking)
		dropResponses = append(dropResponses, dropResponse)
	}

	return wsConn.conn.WriteJSON(dropResponses)
}

func NewDropAvailable(userID uint) {
	fmt.Println("New drop available for user", userID)
	ds := &dropservice.DropService{
		Repo: repositories.Setup(),
	}
	err := sendToUser(userID, ds)
	if err != nil {
		log.Printf("Error sending message to user %d: %v", userID, err)
	}
}

// DeleteDrop godoc
//
//	@Summary		Delete a drop
//	@Description	Delete a drop
//	@Tags			drop
//	@Accept			json
//	@Produce		json
//	@Security BearerAuth
//	@Param			id path int true "Drop ID"
//	@Success		204 No Content
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/drops/:id [delete]
func DeleteDrop(c *gin.Context) {
	repo := repositories.Setup()

	ds := &dropservice.DropService{
		Repo: repo,
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

	id := c.Param("id")

	if "" == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	dropId, err := converters.StringToUint(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ds.DeleteDrop(dropId, uintCurrentUserId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// PatchDrop godoc
//
//	@Summary		Patch a drop
//	@Description	Patch a drop
//	@Tags			drop
//	@Accept			json
//	@Produce		json
//	@Security BearerAuth
//	@Param			id path int true "Drop ID"
//	@Param			drop body		model.DropPatch true "Drop object"
//	@Success		200 {object} response_models.GetDropResponse
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/drops/:id [patch]
func PatchDrop(c *gin.Context) {
	repo := repositories.Setup()

	ds := &dropservice.DropService{
		Repo: repo,
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

	id := c.Param("id")

	if "" == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	dropId, err := converters.StringToUint(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dropPatch model.DropPatch

	if err := c.ShouldBindJSON(&dropPatch); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	updatedDrop, err := ds.PatchDrop(dropId, uintCurrentUserId, dropPatch)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := response_models.FormatGetDropResponse(updatedDrop, false)

	c.JSON(http.StatusOK, response)
}

var hasUserDroppedTodayConnections = make(map[string]*WebSocketConnection)

func HasUserDroppedTodayWS(c *gin.Context) {
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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade WebSocket"})
		return
	}

	wsConn := &WebSocketConnection{conn: conn}

	mu.Lock()
	hasUserDroppedTodayConnections[strconv.Itoa(int(uintCurrentUserId))] = wsConn
	fmt.Printf("Users connected to has user dropped today: %v\n", hasUserDroppedTodayConnections)
	mu.Unlock()

	ds := &dropservice.DropService{
		Repo: repositories.Setup(),
	}

	hasDropped, err := ds.HasUserDroppedToday(uintCurrentUserId)

	if err != nil {
		log.Printf("Error checking if user has dropped today: %v", err)
		return
	}

	err = wsConn.conn.WriteJSON(response_models.HasUserDroppedTodayResponse{Status: hasDropped})

	if err != nil {
		log.Printf("Error sending message to user %d: %v", uintCurrentUserId, err)
		return
	}

	defer func() {
		mu.Lock()
		delete(hasUserDroppedTodayConnections, strconv.Itoa(int(uintCurrentUserId)))
		mu.Unlock()
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

func RefreshHasUserDroppedToday() {
	for {
		mu.Lock()
		for userId := range hasUserDroppedTodayConnections {
			err := hasUserDroppedTodayConnections[userId].conn.WriteJSON(response_models.HasUserDroppedTodayResponse{Status: false})
			if err != nil {
				log.Printf("Error sending message to user %s: %v", userId, err)
				continue
			}
		}
		mu.Unlock()
	}
}

func SearchContentForCurrentDrop(c *gin.Context) {

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

	hadUserDroppedToday, err := ds.HasUserDroppedToday(uintCurrentUserId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if hadUserDroppedToday {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already dropped"})
		return
	}

	search := c.Query("search")

	if "" == search {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search is required"})
		return
	}

	lastDropNotif, err := ds.Repo.DropNotificationRepository.GetCurrentDropNotification()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if nil == lastDropNotif {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No drop notification found"})
		return
	}

	apiService := drop_type_apis.Factory(lastDropNotif.GetType())

	if nil == apiService {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid drop type"})
		return
	}

	apiService.Init()

	results := apiService.Search(search)

	c.JSON(http.StatusOK, results)

}
