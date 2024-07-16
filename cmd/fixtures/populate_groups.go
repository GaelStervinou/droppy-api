package fixtures

import (
	"database/sql"
	"github.com/go-faker/faker/v4"
	"go-api/internal/storage/postgres"
	"gorm.io/gorm"
	"math/rand"
)

func PopulateGroups(db *gorm.DB) {
	activeUsers := make([]postgres.User, 0)
	db.Where("status = ?", 1).Find(&activeUsers)

	images := []string{"https://www.booska-p.com/wp-content/uploads/2022/02/des-mangas-pirate%CC%81s-par-une-firme-americaine-news-visu-1024x750.jpg", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTM1eLN0XU2iTWwxwj1HzO5crOgBeZiTA4C7Q&s", "https://lempreintedigitale.com/wp-content/uploads/2022/03/clubs-foot-europeens-plus-suivis-reseaux-sociaux-min.jpeg", "https://f.hellowork.com/edito/sites/3/2021/10/AdobeStock_419775027-2-1-1200x800.jpeg"}
	memberRoles := postgres.GroupMemberRoles()
	for range 500 {
		grp := postgres.Group{
			Name:        faker.Word(),
			Description: faker.Sentence(),
			CreatedByID: activeUsers[rand.Intn(len(activeUsers))].ID,
			IsPrivate:   rand.Intn(2) == 0,
			PicturePath: sql.NullString{String: images[rand.Intn(len(images))], Valid: true},
		}
		db.Create(&grp)

		for i := range activeUsers {
			j := rand.Intn(i + 1)
			activeUsers[i], activeUsers[j] = activeUsers[j], activeUsers[i]
		}
		randNbUsers := rand.Intn(100)
		for j := range randNbUsers {
			activeUser := activeUsers[j]
			db.Create(&postgres.GroupMember{
				GroupID:  grp.ID,
				MemberID: activeUser.ID,
				Status:   1,
				Role:     memberRoles[rand.Intn(len(memberRoles))],
			})

			var dropIds []uint
			db.Model(&postgres.Drop{}).Where("created_by_id = ?", activeUser.ID).Limit(rand.Intn(5)).Pluck("id", &dropIds)
			for _, dropId := range dropIds {
				db.Create(&postgres.GroupDrop{
					GroupID: grp.ID,
					DropID:  dropId,
				})
			}
		}
	}
}
