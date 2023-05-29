package cloudfirestore

import (
	"cloud.google.com/go/firestore"
)

type FuncGetDoc func(ids ...string) *firestore.DocumentRef

type TypedBatchGetter[T any] interface {
	ConvertibleBatchGetter[T, T]
}
