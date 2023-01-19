package firebaseauth

import (
	"context"

	"firebase.google.com/go/auth"
)

type serviceDebug struct {
	sFirebaseAuth Service
	dummyClaims   map[string]interface{}
}

func NewServiceDebug(cFirebaseAuth *auth.Client, dummyClaims map[string]interface{}) Service {
	sFirebaseAuth := NewService(cFirebaseAuth)
	return &serviceDebug{
		sFirebaseAuth,
		dummyClaims,
	}
}

// 認証を行う
func (s *serviceDebug) Authentication(ctx context.Context, ah string) (string, map[string]interface{}, error) {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	if user := getUserByAuthHeader(ah); user != "" {
		return user, s.dummyClaims, nil
	}
	return s.sFirebaseAuth.Authentication(ctx, ah)
}

// カスタムClaimsを設定
func (s *serviceDebug) SetCustomClaims(ctx context.Context, userID string, claims map[string]interface{}) error {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	ah := getAuthHeader(ctx)
	if getUserByAuthHeader(ah) != "" {
		return nil
	}
	return s.sFirebaseAuth.SetCustomClaims(ctx, userID, claims)
}
