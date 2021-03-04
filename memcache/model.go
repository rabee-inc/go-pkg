package memcache

// Datum ... インメモリキャッシュ
type Datum struct {
	Value     interface{}
	ExpiredAt int64
}
