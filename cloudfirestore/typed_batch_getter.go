package cloudfirestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

type FuncGetDoc func(ids ...string) *firestore.DocumentRef

type TypedBatchGetter[T any] interface {
	Add(ids ...string)
	Delete(ids ...string)
	GetMap() map[string]*T
	Get(ids ...string) *T
	Commit(ctx context.Context) error
}
