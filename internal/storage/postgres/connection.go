package postgres

import (
	"go-api/internal/storage/postgres/token"
	"go-api/internal/storage/postgres/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	//TODO à mettre dans des variables d'env
	dsn := "user=postgres host=localhost port=5432 dbname=postgres password=postgres sslmode=disable"
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
	sqlDB.AutoMigrate(&user.User{}, &token.AuthToken{})
}
