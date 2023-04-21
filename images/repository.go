package images

import "context"

type Repository interface {
	UpdateByConvertObjects(
		ctx context.Context,
		key string,
		objects []*Object,
	) error

	UpdateByGenerateURL(
		ctx context.Context,
		key string,
		id string,
		url string,
	) error
}
