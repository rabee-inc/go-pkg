package firebaseauth

import (
	"context"

	"firebase.google.com/go/auth"
)

type serviceDebug struct {
	svc         Service
	dummyClaims map[string]interface{}
}

// Authentication ... 認証を行う
func (s *serviceDebug) Authentication(ctx context.Context, ah string) (string, map[string]interface{}, error) {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	if user := getUserByAuthHeader(ah); user != "" {
		return user, s.dummyClaims, nil
	}
	return s.svc.Authentication(ctx, ah)
}

// SetCustomClaims ... カスタムClaimsを設定
func (s *serviceDebug) SetCustomClaims(ctx context.Context, userID string, claims map[string]interface{}) error {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	ah := getAuthHeader(ctx)
	if getUserByAuthHeader(ah) != "" {
		return nil
	}
	return s.svc.SetCustomClaims(ctx, userID, claims)
}

func (s *serviceDebug) GetEmail(ctx context.Context, userID string) (string, error) {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	ah := getAuthHeader(ctx)
	if user := getUserByAuthHeader(ah); user != "" {
		return "development@rabee.jp", nil
	}
	return s.svc.GetEmail(ctx, userID)
}

func (s *serviceDebug) GetTwitterID(ctx context.Context, userID string) (string, error) {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	ah := getAuthHeader(ctx)
	if getUserByAuthHeader(ah) != "" {
		return "", nil
	}
	return s.svc.GetTwitterID(ctx, userID)
}

func (s *serviceDebug) ExistUser(ctx context.Context, userID string) (bool, error) {
	return s.svc.ExistUser(ctx, userID)
}

func (s *serviceDebug) IsEmailVerified(ctx context.Context, userID string) (bool, error) {
	return s.svc.IsEmailVerified(ctx, userID)
}

func (s *serviceDebug) CreateUser(ctx context.Context, email string, password string, displayName string) (string, error) {
	return s.svc.CreateUser(ctx, email, password, displayName)
}

func (s *serviceDebug) UpdateUser(ctx context.Context, userID string, email *string, password *string, displayName *string) error {
	return s.svc.UpdateUser(ctx, userID, email, password, displayName)
}

func (s *serviceDebug) DeleteUser(ctx context.Context, userID string) error {
	return s.svc.DeleteUser(ctx, userID)
}

func (s *serviceDebug) GeneratePasswordRemindURL(ctx context.Context, userID string, email string, setting *auth.ActionCodeSettings) (string, error) {
	return s.svc.GeneratePasswordRemindURL(ctx, userID, email, setting)
}

// NewServiceDebug ... ServiceDebugを作成する
func NewServiceDebug(cli *auth.Client, dummyClaims map[string]interface{}) Service {
	svc := NewService(cli)
	return &serviceDebug{
		svc:         svc,
		dummyClaims: dummyClaims,
	}
}
