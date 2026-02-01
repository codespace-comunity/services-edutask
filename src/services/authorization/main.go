package main

import (
	"authorization/controllers"
	"authorization/router"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	if os.Getenv("APP_ENV") != "docker" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("error load env: %v", err)
		}
	}

	// setup gin
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.UseRawPath = true

	r.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	})

	// setup database connection
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	// setup redis
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
	})

	if err := rdb.Ping(ctx); err != nil {
		log.Printf("failed connect to redis: %v", err)
	}

	controller := controllers.New(pool, rdb)

	router.New(r.Group("/api"), controller)

	// running server
	log.Printf("server running on port %s", os.Getenv("APPLICATION_PORT"))
	if err := r.Run(fmt.Sprintf("%s%s", ":", os.Getenv("APPLICATION_PORT"))); err != nil {
		log.Println("failed to running server")
	}
}