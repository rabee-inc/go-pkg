package cloudfirestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

type BatchGetter interface {
	Add(docRef *firestore.DocumentRef, dst any)
	Delete(docRef *firestore.DocumentRef)
	Commit(ctx context.Context) error
	Get(path string) any
}
