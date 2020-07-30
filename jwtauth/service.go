package jwtauth

import (
	"context"
)

// Service ... JWT認証の機能を提供する
type Service interface {
	CreateToken(
		ctx context.Context,
		userID string,
		customClaims map[string]interface{}) (string, error)
	Authentication(
		ctx context.Context,
		ah string) (string, map[string]interface{}, error)
}
