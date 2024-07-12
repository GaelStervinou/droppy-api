package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-api/pkg/jwt_helper"
	"net/http"
	"strings"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		if 2 == len(parts) {
			tokenString := parts[1]
			token, err := jwt_helper.VerifyToken(tokenString)

			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse JWT token: " + err.Error()})
				c.Abort()
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				role := claims["role"].(string)
				if role != "admin" {
					c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
