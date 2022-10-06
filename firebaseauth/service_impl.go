package firebaseauth

import (
	"context"

	"firebase.google.com/go/auth"

	"github.com/rabee-inc/go-pkg/log"
)

type service struct {
	cFirebaseAuth *auth.Client
}

func NewService(cFirebaseAuth *auth.Client) Service {
	return &service{cFirebaseAuth}
}

// Authentication ... 認証を行う
func (s *service) Authentication(ctx context.Context, ah string) (string, map[string]interface{}, error) {
	token := getTokenByAuthHeader(ah)
	if token == "" {
		err := log.Warninge(ctx, "token empty error")
		return "", nil, err
	}

	t, err := s.cFirebaseAuth.VerifyIDToken(ctx, token)
	if err != nil {
		log.Warningf(ctx, "verify token error: %s, %s", token, err.Error())
		return "", nil, err
	}
	return t.UID, t.Claims, nil
}

// SetCustomClaims ... カスタムClaimsを設定
func (s *service) SetCustomClaims(ctx context.Context, userID string, claims map[string]interface{}) error {
	err := s.cFirebaseAuth.SetCustomUserClaims(ctx, userID, claims)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}
