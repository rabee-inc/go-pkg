package inmemcache

import (
	"sync"

	"github.com/rabee-inc/go-pkg/timeutil"
)

type Client[T any] struct {
	itemMap map[string]*Item[T]
	mutex   *sync.Mutex
}

func NewClient[T any]() *Client[T] {
	return &Client[T]{
		map[string]*Item[T]{},
		&sync.Mutex{},
	}
}

func (c *Client[T]) GetOrSet(key string, fn func() (T, int, error)) (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := timeutil.NowUnix()
	if item, ok := c.itemMap[key]; ok {
		if item.ExpiredAt > now {
			// キャッシュを返す
			return item.Value, nil
		}
		delete(c.itemMap, key)
	}
	value, expiredSecond, err := fn()
	if err != nil {
		return value, err
	}
	expiredAt := now + timeutil.SecondsToMilliseconds(expiredSecond)
	c.itemMap[key] = &Item[T]{
		Value:     value,
		ExpiredAt: expiredAt,
	}
	return value, nil
}

func (c *Client[T]) GetMultiOrSetMulti(keys []string, fn func([]string) (map[string]T, int, error)) (map[string]T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := timeutil.NowUnix()
	dsts := map[string]T{}
	missKeys := []string{}
	for _, key := range keys {
		if item, ok := c.itemMap[key]; ok {
			if item.ExpiredAt > now {
				// キャッシュを返す
				dsts[key] = item.Value
				continue
			}
			delete(c.itemMap, key)
		}
		missKeys = append(missKeys, key)
	}
	if len(missKeys) == 0 {
		return dsts, nil
	}
	values, expiredSecond, err := fn(missKeys)
	if err != nil {
		return dsts, err
	}
	expiredAt := now + timeutil.SecondsToMilliseconds(expiredSecond)
	for key, value := range values {
		c.itemMap[key] = &Item[T]{
			Value:     value,
			ExpiredAt: expiredAt,
		}
		dsts[key] = value
	}
	return dsts, nil
}

func (c *Client[T]) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.itemMap, key)
}

func (c *Client[T]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.itemMap = map[string]*Item[T]{}
}
