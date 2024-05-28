package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/cmd/fixtures"
	"go-api/internal/storage/postgres"
	"go-api/internal/storage/postgres/user"
)

func PopulateUsers(c *gin.Context) {
	sqlDB, err := postgres.Connect()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	us := user.NewRepo(sqlDB)

	err = fixtures.PopulateUsers(us)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Users populated"})
}
