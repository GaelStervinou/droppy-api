package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-api/authentication"
	"go-api/authentication/provider"
	_ "go-api/docs"
	"go-api/internal/http/controllers"
	"go-api/internal/storage/postgres"
	"go-api/pkg/jwt_helper"
	"log"
	"net/http"
	"strings"
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
			user.GET("/:id", CurrentUserMiddleware(), controllers.GetUserById)
			user.POST("/", controllers.Create)
		}
	}
	r.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err = r.Run(":3000")

	if err != nil {
		log.Fatal(err)
	}
}

func CurrentUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/auth/") {
			c.Next()
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}
		tokenString := parts[1]
		fmt.Println(tokenString)

		token, err := jwt_helper.VerifyToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse JWT token: " + err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := uint(claims["sub"].(float64))

			fmt.Println(userId)
			//TODO peut-être faire une requête pour récupérer le user et passer le user direct dans le context
			c.Set("userId", userId)
		}
		c.Next()
	}
}
