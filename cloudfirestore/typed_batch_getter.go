package cloudfirestore

type TypedBatchGetter[T any] interface {
	ConvertibleBatchGetter[T, T]
}
