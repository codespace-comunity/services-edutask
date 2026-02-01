package controllers

import (
	"authorization/utils"
	"context"
	"fmt"
	"log"
	"time"
)

func (c *controller) LoginUser(ctx context.Context, request *RequestLoginService) (response *ResponseLoginService, err error) {
	log.Println("login in process...")
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		log.Printf("failed to connect pool db!")
		return nil, err
	}
	defer tx.Rollback(ctx)

	var where string
	var whereVal string

	if request.Email != "" {
		where = "email"
		whereVal = request.Email
	} else if request.Name != "" {
		where = "name"
		whereVal = request.Name
	}

	sql := fmt.Sprintf("SELECT id, name, email, password, created_at FROM users WHERE %s=$1", where)

	var user UserType

	if err := tx.QueryRow(ctx, sql, whereVal).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAT); err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateJWT(user.ID, 15 * time.Minute)
	if err != nil {
		log.Printf("failed to generate access token, arguments: %v", err)
		return nil, err
	}

	refreshToken, err := utils.GenerateJWT(user.ID, 7 * 24 * time.Hour)
	if err != nil {
		log.Printf("failed to generate refresh token, arguments: %v", err)
		return nil, err
	}

	userResponse := UserTypeResponse{
		ID: user.ID,
		Name: request.Name,
		Email: request.Email,
		CreatedAT: time.Now(),
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	if err := c.rdb.Set(ctx, user.ID, refreshToken, 7 * 24 * time.Hour).Err(); err != nil {
		log.Printf("failed to save value in redis, arguments: %v", err)
		return nil, err
	}

	return &ResponseLoginService{
		User: userResponse,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}