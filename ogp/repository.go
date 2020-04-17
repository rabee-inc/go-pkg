package ogp

import "context"

// Repository ... リポジトリ
type Repository interface {
	UpdateURL(ctx context.Context, key string, id string, url string) error
}
