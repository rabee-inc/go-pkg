package inmemcache

type Item[T any] struct {
	Value     T
	ExpiredAt int64
}
