package firebaseauth

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type serviceDebug struct {
	sFirebaseAuth Service
	dummyClaims   map[string]any
}

func NewServiceDebug(cFirebaseAuth *auth.Client, dummyClaims map[string]any) Service {
	sFirebaseAuth := NewService(cFirebaseAuth)
	return &serviceDebug{
		sFirebaseAuth,
		dummyClaims,
	}
}

// 認証を行う
func (s *serviceDebug) Authentication(ctx context.Context, ah string) (string, map[string]any, error) {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	if user := getDebugByAuthHeader(ah); user != "" {
		return user, s.dummyClaims, nil
	}
	return s.sFirebaseAuth.Authentication(ctx, ah)
}

// カスタムClaimsを設定
func (s *serviceDebug) SetCustomClaims(ctx context.Context, userID string, claims map[string]any) error {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	ah := getAuthHeader(ctx)
	if getDebugByAuthHeader(ah) != "" {
		return nil
	}
	return s.sFirebaseAuth.SetCustomClaims(ctx, userID, claims)
}
