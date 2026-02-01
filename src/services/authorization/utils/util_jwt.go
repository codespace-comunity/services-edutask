package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(userID, tokenType string, expired time.Duration) (string, error) {

	claims :=  CustomClaims{
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: uuid.NewString(),
			Issuer: "codespace",
			Subject: userID,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expired)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}

var (
	ErrTokenExpired   = errors.New("token expired")
	ErrTokenInvalid   = errors.New("token invalid")
	ErrTokenMalformed = errors.New("token malformed")
)

func ValidateJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrTokenInvalid
			}
			return []byte(os.Getenv("SECRET")), nil
		},
		jwt.WithIssuer("codespace"),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		default:
			return nil, ErrTokenInvalid
		}
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	// VALIDASI MANUAL TAMBAHAN (tanpa cek Type)
	if claims.Subject == "" {
		return nil, ErrTokenInvalid
	}

	if claims.IssuedAt == nil || claims.ExpiresAt == nil {
		return nil, ErrTokenInvalid
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return claims, nil
}

