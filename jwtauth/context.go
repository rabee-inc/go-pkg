package jwtauth

import "context"

type contextKey string

const (
	userIDContextKey contextKey = "jwtauth:user_id"
	claimsContextKey contextKey = "jwtauth:claims"
)

// GetUserID ... JWTAuthのユーザーIDを取得
func GetUserID(ctx context.Context) string {
	if dst := ctx.Value(userIDContextKey); dst != nil {
		return dst.(string)
	}
	return ""
}

// GetClaims ... JWTAuthのClaimsの値を取得
func GetClaims(ctx context.Context) (map[string]interface{}, bool) {
	if dst := ctx.Value(claimsContextKey); dst != nil {
		return dst.(map[string]interface{}), true
	}
	return nil, false
}

func setUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func setClaims(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}
