package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-api/authentication"
	"go-api/authentication/provider"
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

	provider.UseGoogleAuth()

	postgres.AutoMigrate()

	//var wg sync.WaitGroup
	//repo := repositories.Setup(&wg)
	//defer repo.Disconnect()

	authentication.Init()

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
			user.PATCH("/:id", middlewares.CurrentUserMiddleware(true), controllers.PatchUserById)
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
	err = r.Run(":3000")

	if err != nil {
		log.Fatal(err)
	}
}
