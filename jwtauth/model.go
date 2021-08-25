package jwtauth

import "github.com/golang-jwt/jwt"

// Claims ... JWTに仕込むclaims定義
type Claims struct {
	UserID       string                 `json:"user_id"`
	CustomClaims map[string]interface{} `json:"custom_claims"`
	jwt.StandardClaims
}
