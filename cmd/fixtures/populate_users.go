package fixtures

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"go-api/internal/storage/postgres"
	"go-api/pkg/hash"
	"gorm.io/gorm"
	"math/rand"
)

func PopulateUsers(db *gorm.DB) error {
	avatars := []string{"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSuj5htdFy-s5wzTWvk3ZCZ4KlAKsmEyaA6IQ&s", "https://letstryai.com/wp-content/uploads/2023/11/stable-diffusion-avatar-prompt-example-2.jpg", "https://cdn.prod.website-files.com/61554cf069663530fc823d21/6369fe5c04b5b062eeed6515_download-57-min.png", "https://marketplace.canva.com/EAFltPVX5QA/1/0/1600w/canva-cute-cartoon-anime-girl-avatar-ZHBl2NicxII.jpg", "https://mpost.io/wp-content/uploads/image-7-17.jpg"}
	roles := []string{"user", "admin"}

	hashedPassword, err := hash.GenerateFromPassword("Test123!!")
	db.Create(&postgres.User{
		Email:       "louis@gmail.com",
		Password:    hashedPassword,
		Username:    "Louis",
		Role:        "admin",
		Status:      1,
		IsPrivate:   false,
		Bio:         faker.Sentence(),
		Avatar:      avatars[0],
		FirebaseUID: faker.UUIDDigit(),
	})
	db.Create(&postgres.User{
		Email:       "rui@gmail.com",
		Password:    hashedPassword,
		Username:    "Ruix",
		Role:        "admin",
		Status:      1,
		IsPrivate:   false,
		Bio:         faker.Sentence(),
		Avatar:      avatars[0],
		FirebaseUID: faker.UUIDDigit(),
	})
	db.Create(&postgres.User{
		Email:       "gael@gmail.com",
		Password:    hashedPassword,
		Username:    "Gayelz",
		Role:        "admin",
		Status:      1,
		IsPrivate:   false,
		Bio:         faker.Sentence(),
		Avatar:      avatars[0],
		FirebaseUID: faker.UUIDDigit(),
	})

	for i := range 1000 {

		if err != nil {
			fmt.Printf("Error generating password: %v", err)
		}
		var avatar string
		private := i%2 == 0

		if i%10 != 0 {
			avatar = avatars[rand.Intn(len(avatars))]
		}
		role := roles[rand.Intn(len(roles))]
		db.Create(&postgres.User{
			Email:       faker.Email(),
			Password:    hashedPassword,
			Username:    faker.FirstName(),
			Role:        role,
			Status:      1,
			IsPrivate:   private,
			Bio:         faker.Sentence(),
			Avatar:      avatar,
			FirebaseUID: faker.UUIDDigit(),
		})
		i++
	}
	return nil
}
