package auth

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}
