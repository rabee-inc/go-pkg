package memcache

// Datum ... インメモリキャッシュ
type Datum struct {
	Value     any
	ExpiredAt int64
}
