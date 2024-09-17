package firebaseauth

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func NewClient(projectID string) *auth.Client {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		panic(err)
	}
	cFirebaseAuth, err := app.Auth(ctx)
	if err != nil {
		panic(err)
	}
	return cFirebaseAuth
}
