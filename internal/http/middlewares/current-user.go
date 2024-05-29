package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-api/pkg/jwt_helper"
	"net/http"
	"strings"
)

func CurrentUserMiddleware(forceLogin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			if forceLogin {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
				c.Abort()
				return
			}
			c.Next()
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			if forceLogin {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
				c.Abort()
				return
			}
			c.Next()
		}

		fmt.Println(len(parts))
		if 2 == len(parts) {
			tokenString := parts[1]
			token, err := jwt_helper.VerifyToken(tokenString)

			if err != nil {
				if forceLogin {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse JWT token: " + err.Error()})
					c.Abort()
					return
				}
				c.Next()
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userId := uint(claims["sub"].(float64))
				c.Set("userId", userId)
			}
		}

		c.Next()
	}
}
