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

func (s *PushNotificationService) SendDropNotification(fcmTokens []string, dropType string) {
	if len(fcmTokens) == 0 {
		log.Println("No FCM tokens to send notifications to")
		return
	}

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

	s.SendDropNotification(validTokens, dropType)
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

func (s *PushNotificationService) SendNotification(notifType string, fcmTokens []string) error {
	if len(fcmTokens) == 0 {
		log.Println("No FCM tokens to send notifications to")
		return nil
	}

	firebaseRepo, err := firebase.NewRepo()
	if err != nil {
		log.Println("Error getting firebase repo:", err)
		return err
	}

	client, err := firebaseRepo.App.Messaging(context.Background())
	if err != nil {
		log.Println("Error getting messaging client:", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(fcmTokens) == 1 {
		response2, err2 := client.Send(ctx, &messaging.Message{
			Notification: GetNotificationContent(notifType),
			Token:        fcmTokens[0],
		})

		if err2 != nil {
			log.Println("Error sending single message:", err2)
		} else {
			fmt.Println("Single message response:", response2)
		}
	} else {
		response, err := client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
			Notification: GetNotificationContent(notifType),
			Tokens:       fcmTokens,
		})

		if err != nil {
			log.Println("Error sending multicast message:", err)
		} else {
			fmt.Println("Multicast response:", response)
		}
	}

	return nil
}

func GetNotificationContent(notifType string) *messaging.Notification {
	switch notifType {
	case "follow-public":
		return &messaging.Notification{
			Title: "Nouveau follower !",
			Body:  "Vous avez un nouvel ami",
		}
	case "follow-private":
		return &messaging.Notification{
			Title: "Demande d'ami",
			Body:  "Vous avez une nouvelle demande d'ami",
		}
	case "like":
		return &messaging.Notification{
			Title: "Nouveau like !",
			Body:  "Quelqu'un a aimé votre contenu",
		}
	case "comment":
		return &messaging.Notification{
			Title: "Nouveau commentaire !",
			Body:  "Quelqu'un a commenté votre contenu",
		}
	default:
		return &messaging.Notification{
			Title: "Nouveau post",
			Body:  "Un de vos amis a posté un contenu, allez voir !",
		}
	}
}
