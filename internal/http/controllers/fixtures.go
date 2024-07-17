package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api/cmd/fixtures"
	"go-api/internal/storage/postgres"
)

func PopulateAll(c *gin.Context) {
	truncate := c.Query("truncate")
	sqlDB := postgres.Connect()

	if truncate == "true" {
		fixtures.TruncateTables(sqlDB)
	}

	fixtures.PopulateUsers(sqlDB)
	fixtures.PopulateFollows(sqlDB)
	fixtures.PopulateDrops(sqlDB)
	fixtures.PopulateComments(sqlDB)
	fixtures.PopulateGroups(sqlDB)
}

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

func PopulateGroups(c *gin.Context) {
	sqlDB := postgres.Connect()

	fixtures.PopulateGroups(sqlDB)

	c.JSON(200, gin.H{"message": "Groups populated"})
}

func PopulateComments(c *gin.Context) {
	sqlDB := postgres.Connect()

	err := fixtures.PopulateComments(sqlDB)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Comments populated"})
}
