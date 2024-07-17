package fixtures

import (
	"go-api/internal/storage/postgres"
	"gorm.io/gorm"
	"math/rand"
)

func PopulateComments(db *gorm.DB) error {

	contents := []string{
		"J'adore ce drop!",
		"Je suis fan de Crusty le Clown",
		"M'enfin... jui l'parrain d'ton fils !",
		"J'ai pas de fric, mais j'ai des amis",
		"Je suis un grand fan de la série",
		"Perso je préfère Jul",
		"Le drop est CLOWNESQUE",
	}
	lenContents := len(contents)

	var drops []*postgres.Drop
	db.Where("status = 1").Find(&drops)

	var users []*postgres.User
	db.Where("status = 1").Find(&users)
	lenUsers := len(users)
	for _, drop := range drops {
		for i := 0; i < rand.Intn(10); i++ {
			user := users[rand.Intn(lenUsers)]
			db.Create(&postgres.Comment{
				Content:     contents[rand.Intn(lenContents)],
				CreatedById: user.GetID(),
				DropId:      drop.ID,
			})
		}
	}

	return nil
}
