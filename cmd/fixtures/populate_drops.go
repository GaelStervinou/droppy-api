package fixtures

import (
	"github.com/go-faker/faker/v4"
	"go-api/internal/storage/postgres"
	"go-api/pkg/drop_type_apis"
	"gorm.io/gorm"
	"math/rand/v2"
)

func PopulateDrops(db *gorm.DB) error {
	types := drop_type_apis.GetValidDropTypes()
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
	images := make(map[string][]string)
	images["youtube"] = []string{"https://i.ytimg.com/vi/RLyxAGHGjfg/hqdefault.jpg", "https://www.journaldemickey.com/wp-content/uploads/2023/01/michou.jpg", "https://i.ytimg.com/vi/AI6uPdYDxvo/sddefault.jpg?sqp=-oaymwEmCIAFEOAD8quKqQMa8AEB-AH-CYAC0AWKAgwIABABGGUgXihGMA8=&rs=AOn4CLC8PIKS6ucM5d4LMPhjaCdO2aXWrw"}
	images["twitch"] = []string{"https://static.wikia.nocookie.net/youtuberfrancais/images/6/61/Alderiate.PNG/revision/latest?cb=20200705211516&path-prefix=fr", "https://actustream.fr/img/zen-emission-twitch1.jpg"}
	images["films"] = []string{"https://m.media-amazon.com/images/I/71XlZvKMwoL._AC_UF1000,1000_QL80_.jpg", "https://resize-europe1.lanmedia.fr/img/var/europe1/storage/images/media/images/intouchables/15796690-1-fre-FR/Intouchables_reference.jpg", "https://antreducinema.fr/wp-content/uploads/2020/04/Titanic.jpg", "https://media.gqmagazine.fr/photos/608297ace24bc2c55a7e1c2f/1:1/w_538,h_538,c_limit/plus%20belles%20affiches%20cin%C3%A9ma.png"}
	images["spotify"] = []string{"https://static.fnac-static.com/multimedia/FR/Images_Produits/FR/fnac.com/Visual_Principal_340/9/7/6/3700187626679/tsp20120926064208/Temps-mort.jpg", "https://www.planetegrandesecoles.com/wp-content/uploads/2023/03/jul-parcours-fortune-musique-.png", "https://hips.hearstapps.com/hmg-prod/images/beyonc-c3-a9-performs-onstage-during-the-renaissance-world-news-photo-1707759399.jpg?crop=0.520xw:0.758xh;0.157xw,0&resize=640:*"}
	randomPics := []string{"https://cdn.pixabay.com/photo/2024/05/26/10/15/bird-8788491_1280.jpg", "https://img.freepik.com/photos-gratuite/prise-vue-au-grand-angle-seul-arbre-poussant-sous-ciel-assombri-pendant-coucher-soleil-entoure-herbe_181624-22807.jpg", "https://hips.hearstapps.com/hmg-prod/images/nature-quotes-landscape-1648265299.jpg", "https://upload.wikimedia.org/wikipedia/commons/c/c5/Ben_david.jpg", "https://media.ouest-france.fr/v1/pictures/MjAyNDA1Y2VjYTk2ZDJjYjM3ZGIxYjRmOGY0OWIzNzA1MDQxNzE?width=1260&height=708&focuspoint=50%2C25&cropresize=1&client_id=bpeditorial&sign=06625b8b06b2381f10ccc1bb1ffeb88668b3535a93783611aa475b27bd85a83a"}
	location := []string{"Paris", "Marseille", "Lyon", "Toulouse", "Bordeaux", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen", "Lille", "Nantes", "Rennes", "Strasbourg", "Montpellier", "Grenoble", "Saint-Etienne", "Nice", "Le Havre", "Amiens", "Reims", "Rouen"}

	contents := make(map[string][]string)
	contents["youtube"] = []string{"https://www.youtube.com/watch?v=RLyxAGHGjfg", "https://www.youtube.com/watch?v=AI6uPdYDxvo", "https://www.youtube.com/watch?v=4-8-0-4-0", "https://www.youtube.com/watch?v=RLyxAGHGjfg", "https://www.youtube.com/watch?v=AI6uPdYDxvo", "https://www.youtube.com/watch?v=4-8-0-4-0"}
	contents["twitch"] = []string{"https://www.twitch.tv/videos/123456789", "https://www.twitch.tv/videos/987654321", "https://www.twitch.tv/videos/111111111", "https://www.twitch.tv/videos/123456789", "https://www.twitch.tv/videos/987654321", "https://www.twitch.tv/videos/111111111"}
	contents["films"] = []string{"https://www.youtube.com/watch?v=RLyxAGHGjfg", "https://www.youtube.com/watch?v=AI6uPdYDxvo", "https://www.youtube.com/watch?v=4-8-0-4-0", "https://www.youtube.com/watch?v=RLyxAGHGjfg", "https://www.youtube.com/watch?v=AI6uPdYDxvo", "https://www.youtube.com/watch?v=4-8-0-4-0"}
	contents["spotify"] = []string{
		"https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC?si=8b1b1b1b1b1b1b1b",
		"https://open.spotify.com/track/15RB3lFt2Mhc16m5fTTYkh?si=c556aeb639ea4d22",
		"https://open.spotify.com/track/6woLe6dqAdDlA3yxmkR4EO?si=b2c616f232744abd",
	}
	for i := range dropNotifications {
		iuint := uint64(uint(i))
		for j := range activeUsers {
			s := rand.NewPCG(iuint, iuint*47329*uint64(j))
			r := rand.New(s)
			isPinned := r.IntN(3) == 0
			dropType := dropNotifications[i].Type
			picture := images[dropType][r.IntN(len(images[dropType]))]
			randomPic := randomPics[r.IntN(len(randomPics))]
			content := contents[dropType][r.IntN(len(contents[dropType]))]
			db.Create(&postgres.Drop{
				Type:               dropType,
				ContentTitle:       faker.Sentence(),
				ContentSubtitle:    faker.Sentence(),
				Content:            content,
				ContentPicturePath: picture,
				Description:        faker.Sentence(),
				CreatedById:        activeUsers[j].ID,
				Status:             1,
				IsPinned:           isPinned,
				DropNotificationID: dropNotifications[i].ID,
				Lat:                faker.Latitude(),
				Lng:                faker.Longitude(),
				Location:           location[r.IntN(len(location))],
				PicturePath:        randomPic,
			})
		}
	}

	return nil
}
