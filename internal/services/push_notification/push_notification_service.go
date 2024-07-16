package push_notification

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"go-api/internal/repositories"
	"go-api/internal/storage/firebase"
	"go-api/internal/storage/postgres"
	"log"
	"math/rand"
	"time"
)

type PushNotificationService struct {
	Repo *repositories.Repositories
}

func (s *PushNotificationService) SendNotification(fcmTokens []string, dropType string) {
	if len(fcmTokens) == 0 {
		log.Println("No FCM tokens to send notifications to")
		return
	}

	fmt.Println(len(fcmTokens), "FCM tokens found")
	fmt.Println(fcmTokens[0])

	firebaseRepo, err := firebase.NewRepo()
	if err != nil {
		log.Println("Error getting firebase repo:", err)
		return
	}

	client, err := firebaseRepo.App.Messaging(context.Background())
	if err != nil {
		log.Println("Error getting messaging client:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(fcmTokens) == 1 {
		response2, err2 := client.Send(ctx, &messaging.Message{
			Notification: &messaging.Notification{
				Title: "Nouveau Drop !",
				Body:  "Postez vite votre contenu " + dropType,
			},
			Token: fcmTokens[0],
		})

		if err2 != nil {
			log.Println("Error sending single message:", err2)
		} else {
			fmt.Println("Single message response:", response2)
		}
	} else {
		response, err := client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
			Notification: &messaging.Notification{
				Title: "Nouveau Drop !",
				Body:  "Postez vite votre contenu " + dropType,
			},
			Tokens: fcmTokens,
		})

		if err != nil {
			log.Println("Error sending multicast message:", err)
		} else {
			fmt.Println("Multicast response:", response)
		}
	}

	fmt.Println("FCM tokens used:", fcmTokens)
}

func (s *PushNotificationService) SendNotificationsToAllUser(dropType string) {
	sqlDB := postgres.Connect()

	userRepo := postgres.NewUserRepo(sqlDB)
	fcmTokens, err := userRepo.GetAllFCMTokens()
	if err != nil {
		log.Println("Error getting FCM tokens:", err)
		return
	}

	// Remove empty FCM tokens
	var validTokens []string
	for _, token := range fcmTokens {
		if token != "" {
			validTokens = append(validTokens, token)
		}
	}

	s.SendNotification(validTokens, dropType)
}

func (s *PushNotificationService) GenerateRandomNotification(dropType string) {
	startHour := 8
	endHour := 21

	// Generate a random hour and minute within the specified range
	randomHour := rand.Intn(endHour-startHour) + startHour
	randomMinute := rand.Intn(60)

	// Calculate the duration until the random time
	now := time.Now()
	randomTime := time.Date(now.Year(), now.Month(), now.Day(), randomHour, randomMinute, 0, 0, now.Location())
	if randomTime.Before(now) {
		randomTime = randomTime.Add(24 * time.Hour) // Schedule for the next day if the time has already passed
	}

	duration := randomTime.Sub(now)
	fmt.Printf("Notification scheduled for: %v\n", randomTime)

	// Sleep until the random time and then send the notification
	time.Sleep(duration)

	s.SendNotificationsToAllUser(dropType)
}
