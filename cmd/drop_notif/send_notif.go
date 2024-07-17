package drop_notif

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/google/martian/v3/log"
	"go-api/internal/storage/firebase"
	"go-api/internal/storage/postgres"
	"math/rand"
	"time"
)

func GenerateRandomNotification() {
	startHour := 8
	endHour := 21

	// Generate a random hour and minute within the specified range
	rand.New(rand.NewSource(time.Now().UnixNano()))
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

	SendNotificationsToAllUser()
}

func SendNotificationsToAllUser() {
	userRepo := postgres.NewUserRepo(postgres.Connect())
	fcmTokens, err := userRepo.GetAllFCMTokens()
	if err != nil {
		log.Errorf("error getting fcm tokens: %v", err)
	}

	SendNotification(fcmTokens)
}

func SendNotification(fcmTokens []string) {
	firebaseRepo, err := firebase.NewRepo()

	if err != nil {
		panic(err)
	}

	client, err := firebaseRepo.App.Messaging(context.Background())

	if err != nil {
		log.Errorf("error getting messaging client: %v", err)
	}

	response, err := client.SendEachForMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: "Congratulations!!",
			Body:  "You have just implemented push notification",
		},
		Tokens: fcmTokens,
	})

	if err != nil {
		log.Errorf("error sending messages from firebase: %v", err)
	}

	log.Debugf("Successfully sent notification: %v", response)
}
