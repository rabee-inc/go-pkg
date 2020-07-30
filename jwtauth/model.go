package jwtauth

import jwt "github.com/dgrijalva/jwt-go"

// Claims ... JWTに仕込むclaims定義
type Claims struct {
	UserID       string                 `json:"user_id"`
	CustomClaims map[string]interface{} `json:"custom_claims"`
	jwt.StandardClaims
}
