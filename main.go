package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log2 "github.com/google/martian/v3/log"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-api/cmd/drop_notif"
	_ "go-api/docs"
	"go-api/internal/http/controllers"
	"go-api/internal/http/middlewares"
	"go-api/internal/storage/postgres"
	"log"
)

// @title Droppy API
// @version 1.0
// @description This is the API documentation for Droppy

// @contact.name   Droppy API Support
// @contact.email  stervinou.g36@gmail.com
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	postgres.AutoMigrate()
	r := gin.Default()
	r.Use(gin.Recovery())

	v1 := r.Group("/")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/refresh", controllers.RefreshToken)
			auth.GET("/:provider", controllers.GoogleAuth)
			auth.GET("/:provider/callback", controllers.GoogleAuthCallback)
			auth.POST("/", controllers.Login)
			auth.POST("/oauth_token", controllers.FirebaseLogin)
		}

		user := v1.Group("/users")
		{
			user.GET("/:id", middlewares.CurrentUserMiddleware(true), controllers.GetUserById)
			user.POST("/", controllers.Create)
			user.GET("/search", controllers.SearchUsers)
			user.PATCH("/:id", middlewares.CurrentUserMiddleware(true), controllers.PatchUserById)
			user.GET("/my-feed/ws", middlewares.CurrentUserMiddleware(true), controllers.GetCurrentUserFeedWS)
			user.GET("/my-feed", middlewares.CurrentUserMiddleware(true), controllers.GetCurrentUserFeed)
			user.GET("/:id/drops", middlewares.CurrentUserMiddleware(true), controllers.DropsByUserId)

			user.GET("/:id/following", middlewares.CurrentUserMiddleware(true), controllers.GetUserFollowing)
			user.GET("/:id/followers", middlewares.CurrentUserMiddleware(true), controllers.GetUserFollowers)
		}

		follow := v1.Group("/follows")
		{
			follow.POST("/", middlewares.CurrentUserMiddleware(true), controllers.FollowUser)
			follow.GET("/pending", middlewares.CurrentUserMiddleware(true), controllers.GetMyPendingRequestsWS)
			follow.POST("/accept/:id", middlewares.CurrentUserMiddleware(true), controllers.AcceptRequest)
			follow.POST("/reject/:id", middlewares.CurrentUserMiddleware(true), controllers.RejectRequest)
			//follow.GET("/:id/accept", middlewares.CurrentUserMiddleware(), controllers.AcceptFollow)
		}

		drop := v1.Group("/drops")
		{
			drop.POST("/", middlewares.CurrentUserMiddleware(true), controllers.CreateDrop)
			drop.POST("/:id/comments", middlewares.CurrentUserMiddleware(true), controllers.CommentDrop)
			drop.POST("/:id/like", middlewares.CurrentUserMiddleware(true), controllers.LikeDrop)
			drop.DELETE("/:id/like", middlewares.CurrentUserMiddleware(true), controllers.UnlikeDrop)
		}

		group := v1.Group("/groups")
		{
			group.POST("/", middlewares.CurrentUserMiddleware(true), controllers.CreateGroup)
			group.GET("/:id", middlewares.CurrentUserMiddleware(true), controllers.GetOneGroup)
			group.GET("/:id/feed", middlewares.CurrentUserMiddleware(true), controllers.GetGroupFeed)
			group.PATCH("/:id", middlewares.CurrentUserMiddleware(true), controllers.PatchGroup)
			group.GET("/search", middlewares.CurrentUserMiddleware(true), controllers.SearchGroups)

			group.POST("/members/:id/join", middlewares.CurrentUserMiddleware(true), controllers.JoinGroup)
			group.POST("/members/:id/:userId", middlewares.CurrentUserMiddleware(true), controllers.AddUserToGroup)
			group.PATCH("/members/:groupId/:memberId", middlewares.CurrentUserMiddleware(true), controllers.PatchGroupMember)
			group.DELETE("/members/:groupId/:memberId", middlewares.CurrentUserMiddleware(true), controllers.DeleteGroupMember)
		}

		comment := v1.Group("/comments")
		{
			comment.POST("/:id/responses", middlewares.CurrentUserMiddleware(true), controllers.RespondToComment)
			comment.DELETE("/:id/responses/:responseId", middlewares.CurrentUserMiddleware(true), controllers.DeleteCommentResponse)
		}

		fixtures := v1.Group("/fixtures")
		{
			fixtures.GET("/users", controllers.PopulateUsers)
			fixtures.GET("/follows", controllers.PopulateFollows)
			fixtures.GET("/drops", controllers.PopulateDrops)
		}

		admin := v1.Group("/admin")
		{
			admin.GET("/users", middlewares.AdminRequired(), controllers.GetAllUsers)
			admin.GET("/groups", middlewares.AdminRequired(), controllers.GetAllGroups)
			//admin.GET("/drops", controllers.GetAllDrops)
			//admin.GET("/follows", controllers.GetAllFollows)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	c := cron.New()
	_, err = c.AddFunc("0 0 * * *", drop_notif.GenerateRandomNotification)

	if err != nil {
		log2.Errorf("Error adding cron job: %v", err)
	}

	c.Start()
	fmt.Println("Scheduler started...")

	err = r.Run(":3000")

	if err != nil {
		log.Fatal(err)
	}
}
