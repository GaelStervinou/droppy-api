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
// @BasePath /api/v1
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

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.GET("/refresh", controllers.RefreshToken)
			auth.GET("/:provider", controllers.GoogleAuth)
			auth.GET("/:provider/callback", controllers.GoogleAuthCallback)
			auth.POST("/login", controllers.Login)
		}

		user := v1.Group("/users")
		{
			user.GET("/:id", middlewares.CurrentUserMiddleware(), controllers.GetUserById)
			user.POST("/", controllers.Create)
			user.PATCH("/:id", middlewares.CurrentUserMiddleware(), controllers.PatchUserById)
		}
	}
	r.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err = r.Run(":3000")

	if err != nil {
		log.Fatal(err)
	}
}
