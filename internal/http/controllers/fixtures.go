package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/cmd/fixtures"
	"go-api/internal/storage/postgres"
)

func PopulateUsers(c *gin.Context) {
	sqlDB := postgres.Connect()

	err := fixtures.PopulateUsers(sqlDB)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Users populated"})
}

func PopulateFollows(c *gin.Context) {
	sqlDB := postgres.Connect()

	err := fixtures.PopulateFollows(sqlDB)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Follows populated"})
}

func PopulateDrops(c *gin.Context) {
	sqlDB := postgres.Connect()

	err := fixtures.PopulateDrops(sqlDB)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Drops populated"})
}
