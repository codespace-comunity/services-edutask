package controllers

import (
	"authorization/utils"
	"context"
	"errors"
	"log"
	"time"
)

func (c *controller) LoginUser(ctx context.Context, request *RequestLoginService) (response *ResponseLoginService, err error) {
	log.Println("login in process...")

	var (
		sql string
		val string
	)

	if request.Email != "" {
		sql = "SELECT id, name, email, password, created_at FROM users WHERE email=$1"
		val = request.Email
	} else if request.Name != "" {
		sql = "SELECT id, name, email, password, created_at FROM users WHERE name=$1"
		val = request.Name
	}

	var user UserType

	if err := c.pool.QueryRow(ctx, sql, val).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAT); err != nil {
		return nil, err
	}

	if !utils.CompareStringWithHash(request.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	accessToken, err := utils.GenerateJWT(user.ID, "access_token", 15 * time.Minute)
	if err != nil {
		log.Printf("failed to generate access token, arguments: %v", err)
		return nil, err
	}

	refreshToken, err := utils.GenerateJWT(user.ID, "refresh_token", 7 * 24 * time.Hour)
	if err != nil {
		log.Printf("failed to generate refresh token, arguments: %v", err)
		return nil, err
	}

	userResponse := UserTypeResponse{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		CreatedAT: user.CreatedAT,
	}

	c.rdb.Set(ctx, user.ID, refreshToken, 7 * 24 * time.Hour)

	return &ResponseLoginService{
		User: userResponse,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}