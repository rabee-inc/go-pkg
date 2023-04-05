package cloudfirestore

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

func NewClient(projectID string) *firestore.Client {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		panic(err)
	}
	cFirestore, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	return cFirestore
}
