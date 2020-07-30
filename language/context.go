package language

import (
	"context"
)

type contextKey string

const (
	keyContextKey contextKey = "language:key"
)

func getKey(ctx context.Context) Key {
	if dst := ctx.Value(keyContextKey); dst != nil {
		return dst.(Key)
	}
	return ""
}

func setKey(ctx context.Context, key Key) context.Context {
	return context.WithValue(ctx, keyContextKey, key)
}
