package utils

import "github.com/golang-jwt/jwt/v5"

type (
	CustomClaims struct {
		Type string `json:"type"`
		jwt.RegisteredClaims
	}
)