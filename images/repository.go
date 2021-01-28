package images

import "context"

// Repository ... リポジトリ
type Repository interface {
	UpdateByConvertObjects(ctx context.Context, key string, objects []*Object) error
	UpdateByGenerateURL(ctx context.Context, key string, id string, url string) error
}
