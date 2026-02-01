package controllers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type (
	controller struct {
		pool *pgxpool.Pool
		rdb *redis.Client
	}
)

func New(
	pool *pgxpool.Pool,
	rdb *redis.Client,
) IController {
	return &controller{
		pool: pool,
		rdb: rdb,
	}
}