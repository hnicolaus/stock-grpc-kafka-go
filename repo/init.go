package repo

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type Repo struct {
	redisClient RedisClient
}

func New() *Repo {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password if not set
		DB:       0,                // Use default DB
	})

	// Ping the Redis server to check if it's reachable
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Print("Failed to connect to Redis:", err)
		return nil
	}

	return &Repo{
		redisClient: client,
	}
}
