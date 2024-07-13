package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DB *gorm.DB

func Init() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dsn := "user=" + dbUser + " host=" + dbHost + " dbname=" + dbName + " password=" + dbPassword + " sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func Connect() *gorm.DB {
	return DB
}

func AutoMigrate() {
	sqlDB := Connect()

	sqlDB.AutoMigrate(
		&User{},
		&AuthToken{},
		&Follow{},
		&Drop{},
		&DropNotification{},
		&Group{},
		&GroupMember{},
		&GroupDrop{},
		&Comment{},
		&CommentResponse{},
		&Like{},
	)
	fmt.Println("Migrations done")
}
