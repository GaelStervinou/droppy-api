package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
)

func NewRepo() (*FirebaseRepo, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &FirebaseRepo{
		App: app,
	}, nil
}

type FirebaseRepo struct {
	App *firebase.App
}
