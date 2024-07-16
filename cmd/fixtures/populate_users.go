package fixtures

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"go-api/internal/storage/postgres"
	"go-api/pkg/hash"
	"gorm.io/gorm"
	"math/rand/v2"
)

func PopulateUsers(db *gorm.DB) error {
	statuses := []int{-2, -1, 0, 1}
	avatars := []string{"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSuj5htdFy-s5wzTWvk3ZCZ4KlAKsmEyaA6IQ&s", "https://letstryai.com/wp-content/uploads/2023/11/stable-diffusion-avatar-prompt-example-2.jpg", "https://cdn.prod.website-files.com/61554cf069663530fc823d21/6369fe5c04b5b062eeed6515_download-57-min.png", "https://marketplace.canva.com/EAFltPVX5QA/1/0/1600w/canva-cute-cartoon-anime-girl-avatar-ZHBl2NicxII.jpg", "https://mpost.io/wp-content/uploads/image-7-17.jpg"}
	roles := []string{"user", "admin"}

	for i := range 1000 {
		hashedPassword, err := hash.GenerateFromPassword("Test123!!")

		if err != nil {
			fmt.Printf("Error generating password: %v", err)
		}
		iuint := uint64(uint(i))
		s := rand.NewPCG(iuint, iuint*47329)
		r := rand.New(s)
		var avatar string
		private := i%2 == 0
		// initialize local pseudorandom generator
		status := statuses[r.IntN(len(statuses))]
		if i%10 != 0 {
			avatar = avatars[r.IntN(len(avatars))]
		}
		role := roles[r.IntN(len(roles))]
		db.Create(&postgres.User{
			Email:       faker.Email(),
			Password:    hashedPassword,
			Username:    faker.FirstName(),
			Role:        role,
			Status:      status,
			IsPrivate:   private,
			Bio:         faker.Sentence(),
			Avatar:      avatar,
			FirebaseUID: faker.UUIDDigit(),
		})
		i++
	}
	return nil
}
