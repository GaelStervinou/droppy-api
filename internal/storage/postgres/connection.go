package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func Connect() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dsn := "user=" + dbUser + " host=" + dbHost + " dbname=" + dbName + " password=" + dbPassword + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AutoMigrate() {
	sqlDB, err := Connect()
	if err != nil {
		panic(err)
	}
	sqlDB.AutoMigrate(
		&User{},
		&AuthToken{},
		&Follow{},
		&Drop{},
		&DropNotification{},
		&Group{},
		&GroupMember{},
		&Comment{},
		&CommentResponse{},
		&Like{},
	)
	fmt.Println("Migrations done")
}
