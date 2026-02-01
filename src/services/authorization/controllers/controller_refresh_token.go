package controllers

import (
	"authorization/utils"
	"context"
	"errors"
	"log"
	"time"
)

func (c *controller) RefreshToken(ctx context.Context, request *RequestRefreshTokenService) (response *ResponseRefreshTokenService, err error) {
	claims, err := utils.ValidateJWT(request.RefreshToken)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh_token" {
		return nil, errors.New("invalid token")
	}

	storedToken, err := c.rdb.Get(ctx, claims.Subject).Result()
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if storedToken != request.RefreshToken {
		return nil, errors.New("invalid token")
	}

	accessToken, err := utils.GenerateJWT(claims.Subject, "access_token", 15 * time.Minute)
	if err != nil {
		log.Printf("failed to generate access token, arguments: %v", err)
		return nil, err
	}

	return &ResponseRefreshTokenService{
		AccessToken: accessToken,
	}, nil
}