package language

import (
	"context"
)

// GetKey ... 種類を取得する
func GetKey(ctx context.Context) Key {
	return getKey(ctx)
}

// SetKey ... 種類を設定する
func SetKey(ctx context.Context, key Key) context.Context {
	return setKey(ctx, key)
}

// Get ... 現在の言語環境の言語を取得する
func Get(ctx context.Context, key string, texts map[string]Text) string {
	if text, ok := texts[key]; ok {
		key := getKey(ctx)
		if dst, ok := text[key]; ok {
			return dst
		}
	}
	return key
}
