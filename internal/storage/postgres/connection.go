package postgres

import (
	"go-api/internal/storage/postgres/drop"
	"go-api/internal/storage/postgres/drop_notification"
	"go-api/internal/storage/postgres/follow"
	"go-api/internal/storage/postgres/token"
	"go-api/internal/storage/postgres/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func Connect() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dsn := "user=" + dbUser + " host=" + dbHost + " dbname=" + dbName + " password=" + dbPassword + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
	sqlDB.AutoMigrate(&user.User{}, &token.AuthToken{}, &follow.Follow{}, &drop.Drop{}, &drop_notification.DropNotification{})
}
