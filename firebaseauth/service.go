package firebaseauth

import (
	"context"
)

// Service ... Firebase認証の機能を提供する
type Service interface {
	Authentication(
		ctx context.Context,
		ah string) (string, map[string]interface{}, error)
	SetCustomClaims(
		ctx context.Context,
		userID string,
		claims map[string]interface{}) error
	GetEmail(
		ctx context.Context,
		userID string) (string, error)
	GetTwitterID(
		ctx context.Context,
		userID string) (string, error)
	ExistUser(
		ctx context.Context,
		userID string) (bool, error)
	CreateUser(
		ctx context.Context,
		email string,
		password string,
		displayName string) (string, error)
	UpdateUser(
		ctx context.Context,
		userID string,
		email *string,
		password *string,
		displayName *string) error
	DeleteUser(
		ctx context.Context,
		userID string) error
}
