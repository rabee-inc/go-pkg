package jwtauth

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/rabee-inc/go-pkg/errcode"
	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/stringutil"
)

type service struct {
	signKey []byte
	expired time.Duration
}

// CreateToken ... トークンを作成する
func (s *service) CreateToken(ctx context.Context, userID string, customClaims map[string]interface{}) (string, error) {
	// 有効期限
	now := time.Now()

	// トークンを作成
	claims := &Claims{
		UserID:       userID,
		CustomClaims: customClaims,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(s.expired).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 署名
	signedToken, err := token.SignedString(s.signKey)
	if err != nil {
		log.Error(ctx, err)
		return "", err
	}
	return signedToken, nil
}

// Authentication ... 認証を行う
func (s *service) Authentication(ctx context.Context, ah string) (string, map[string]interface{}, error) {
	// 署名済みトークンを取得
	signedToken := getTokenByAuthHeader(ah)
	if signedToken == "" {
		err := log.Warninge(ctx, "empty token")
		err = errcode.Set(err, http.StatusForbidden)
		return "", nil, err
	}

	// トークンを取得
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.signKey, nil
	})
	if err != nil {
		log.Warning(ctx, err)
		return "", nil, err
	}

	// 有効性を確認
	if !token.Valid {
		err := log.Warningc(ctx, http.StatusForbidden, "invalid token")
		return "", nil, err
	}

	// 有効期限を確認
	now := time.Now().Unix()
	if claims.ExpiresAt < now {
		err := log.Warningc(ctx, http.StatusForbidden, "expired token")
		return "", nil, err
	}
	return claims.UserID, claims.CustomClaims, nil
}

// NewService ... Serviceを作成する
func NewService(signKey string, expired time.Duration) Service {
	bSignKey := stringutil.ToBytes(signKey)
	return &service{
		signKey: bSignKey,
		expired: expired,
	}
}
