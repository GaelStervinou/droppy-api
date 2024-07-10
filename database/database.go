package database

import (
	postgres2 "go-api/internal/storage/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	// Connect to the database with gorm
	// return the connection
	dsn := "user=postgres host=localhost port=5432 dbname=postgres password=pass sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&postgres2.User{})

	return db
}
