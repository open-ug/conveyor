/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package api

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost" + ":" + os.Getenv("CRANE_REDIS_PORT"),
		Password: os.Getenv("CRANE_REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
	return client
}
