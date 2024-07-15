package fixtures

import (
	"github.com/go-faker/faker/v4"
	"go-api/internal/storage/postgres"
	"gorm.io/gorm"
	"math/rand/v2"
)

func PopulateDrops(db *gorm.DB) error {

	types := []string{"youtube", "twitch", "films", "spotify", "tiktok"}
	var dropNotifications []postgres.DropNotification
	for i := range 365 {
		iuint := uint64(uint(i))
		s := rand.NewPCG(iuint, iuint*47329)
		r := rand.New(s)
		dropNotifications = append(dropNotifications, postgres.DropNotification{
			Type: types[r.IntN(len(types))],
		})
		db.Create(&dropNotifications[i])
	}

	var activeUsers []postgres.User
	db.Where("status = ?", 1).Find(&activeUsers)
	statuses := []int{-1, 1}
	images := make(map[string][]string)
	images["youtube"] = []string{"https://i.ytimg.com/vi/RLyxAGHGjfg/hqdefault.jpg", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTPuJxaPIWJT4qdma0OwtXhJwS6dy--hYD6ab52Et0RNG29qwepYVPKR1kIKxgNR4ibDic&usqp=CAU", "https://i.ytimg.com/vi/AI6uPdYDxvo/sddefault.jpg?sqp=-oaymwEmCIAFEOAD8quKqQMa8AEB-AH-CYAC0AWKAgwIABABGGUgXihGMA8=&rs=AOn4CLC8PIKS6ucM5d4LMPhjaCdO2aXWrw"}
	images["twitch"] = []string{"https://c.clc2l.com/t/t/w/twitch-4aRVhk.png", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTrkPx4Ow5yqjiymd4_5v-Y92jYeVtKtUqvJw&s"}
	images["films"] = []string{"https://img.freepik.com/photos-gratuite/vue-salle-cinema-3d_23-2151067056.jpg"}
	images["tiktok"] = []string{"https://pbs.twimg.com/media/F935U36XwAAqm46?format=jpg&name=large", "https://img.20mn.fr/ST4yso4CS2meAkgSOFsB7yk/1444x920_des-extraits-de-regards-echanges-par-jordan-bardella-et-gabriel-attal-sont-utilises-pour-faire-des-montages-videos-mettant-en-scene-une-romance-entre-eux"}
	images["spotify"] = []string{"https://play-lh.googleusercontent.com/Gk-KGYaWDqWnAY8UdsmJIqtai3lPBo0CGO20plP43B0VV7ifqr4APihwWVHcLhJCoyfE", "https://www.planetegrandesecoles.com/wp-content/uploads/2023/03/jul-parcours-fortune-musique-.png"}
	randomPic := "https://picsum.photos/300/200"
	location := []string{"Paris", "Marseille", "Lyon", "Toulouse", "Bordeaux", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen"}
	for i := range dropNotifications {
		iuint := uint64(uint(i))
		for j := range activeUsers {
			s := rand.NewPCG(iuint, iuint*47329*uint64(j))
			r := rand.New(s)
			isPinned := r.IntN(3) == 0
			dropType := dropNotifications[i].Type
			status := statuses[r.IntN(len(statuses))]
			picture := images[dropType][r.IntN(len(images[dropType]))]
			hasPicture := r.IntN(3) == 0
			var dropPicture string
			if hasPicture {
				dropPicture = randomPic
			}
			db.Create(&postgres.Drop{
				Type:               dropType,
				ContentTitle:       faker.Sentence(),
				ContentSubtitle:    faker.Sentence(),
				Content:            "content",
				ContentPicturePath: picture,
				Description:        faker.Sentence(),
				CreatedById:        activeUsers[j].ID,
				Status:             uint(status),
				IsPinned:           isPinned,
				DropNotificationID: dropNotifications[i].ID,
				Lat:                faker.Latitude(),
				Lng:                faker.Longitude(),
				Location:           location[r.IntN(len(location))],
				PicturePath:        dropPicture,
			})
		}
	}

	return nil
}
