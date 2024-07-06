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
			user.GET("/:id", middlewares.CurrentUserMiddleware(false), controllers.GetUserById)
			user.POST("/", controllers.Create)
			user.GET("/search", controllers.SearchUsers)
			user.PATCH("/:id", middlewares.CurrentUserMiddleware(true), controllers.PatchUserById)
			user.GET("/my-feed", middlewares.CurrentUserMiddleware(true), controllers.GetCurrentUserFeed)
		}

		follow := v1.Group("/follows")
		{
			follow.POST("/", middlewares.CurrentUserMiddleware(true), controllers.FollowUser)
			//follow.GET("/:id/accept", middlewares.CurrentUserMiddleware(), controllers.AcceptFollow)
		}

		drop := v1.Group("/drops")
		{
			drop.POST("/", middlewares.CurrentUserMiddleware(true), controllers.CreateDrop)
		}

		fixtures := v1.Group("/fixtures")
		{
			fixtures.GET("/users", controllers.PopulateUsers)
			fixtures.GET("/follows", controllers.PopulateFollows)
			fixtures.GET("/drops", controllers.PopulateDrops)
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
