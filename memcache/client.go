package memcache

import (
	"sync"

	"github.com/rabee-inc/go-pkg/timeutil"
)

// Client ... クライアント
type Client struct {
	data  map[string]*Datum
	mutex *sync.Mutex
}

// GetOrSet ... キャッシュを取得（設定）する
func (c *Client) GetOrSet(key string, fn func(key string) (interface{}, int, error)) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := timeutil.NowUnix()
	if datum, ok := c.data[key]; ok {
		if datum.ExpiredAt > now {
			return datum.Value, nil
		}
		delete(c.data, key)
	}
	value, expiredMinute, err := fn(key)
	if err != nil {
		return nil, err
	}
	expiredAt := now + timeutil.MinutesToMilliseconds(expiredMinute)
	datum := &Datum{
		Value:     value,
		ExpiredAt: expiredAt,
	}
	c.data[key] = datum
	return value, nil
}

// NewClient ... クライアントを作成する
func NewClient() *Client {
	data := map[string]*Datum{}
	mutex := &sync.Mutex{}
	return &Client{
		data:  data,
		mutex: mutex,
	}
}
