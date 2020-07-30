package jwtauth

import (
	"context"
	"time"
)

type serviceDebug struct {
	svc         Service
	dummyClaims map[string]interface{}
}

// CreateToken ... トークンを作成する
func (s *serviceDebug) CreateToken(ctx context.Context, userID string, customClaims map[string]interface{}) (string, error) {
	return s.svc.CreateToken(ctx, userID, customClaims)
}

// Authentication ... 認証を行う
func (s *serviceDebug) Authentication(ctx context.Context, ah string) (string, map[string]interface{}, error) {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	if user := getUserByAuthHeader(ah); user != "" {
		return user, s.dummyClaims, nil
	}
	return s.svc.Authentication(ctx, ah)
}

// NewServiceDebug ... Serviceを作成する
func NewServiceDebug(signKey string, expired time.Duration, dummyClaims map[string]interface{}) Service {
	svc := NewService(signKey, expired)
	return &serviceDebug{
		svc:         svc,
		dummyClaims: dummyClaims,
	}
}
