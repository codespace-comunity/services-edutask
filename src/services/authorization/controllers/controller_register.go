package controllers

import (
	"authorization/utils"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

func (c *controller) RegisterUser(ctx context.Context, request *RequestRegisterService) (response *ResponseRegisterService, err error) {
	log.Println("register in process...")
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		log.Printf("failed to connect pool, arguments: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	sql := `INSERT INTO users
			(id, name, email, password, created_at)
			VALUES ($1,$2,$3,$4,$5)`

	hashPassword, err := utils.HashString(request.Password)
	if err != nil {
		log.Printf("failed to hashing password, arguments: %v", err)
		return nil, err
	}

	userID := uuid.NewString()

	_, err = tx.Exec(ctx, sql, userID, request.Name, request.Email, hashPassword, time.Now())
	if err != nil {
		log.Printf("failed to insert data, arguments: %v", err)
		return nil, err
	}

	accessToken, err := utils.GenerateJWT(userID, 15 * time.Minute)
	if err != nil {
		log.Printf("failed to generate access token, arguments: %v", err)
		return nil, err
	}

	refreshToken, err := utils.GenerateJWT(userID, 7 * 24 * time.Hour)
	if err != nil {
		log.Printf("failed to generate refresh token, arguments: %v", err)
		return nil, err
	}

	user := UserTypeResponse{
		ID: userID,
		Name: request.Name,
		Email: request.Email,
		CreatedAT: time.Now(),
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	if err := c.rdb.Set(ctx, userID, refreshToken, 7 * 24 * time.Hour).Err(); err != nil {
		log.Printf("failed to save value in redis, arguments: %v", err)
		return nil, err
	}

	return &ResponseRegisterService{
		User: user,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}