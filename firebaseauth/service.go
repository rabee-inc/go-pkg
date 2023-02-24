package firebaseauth

import (
	"context"
)

type Service interface {
	Authentication(
		ctx context.Context,
		ah string,
	) (string, map[string]any, error)

	SetCustomClaims(
		ctx context.Context,
		userID string,
		claims map[string]any,
	) error
}
