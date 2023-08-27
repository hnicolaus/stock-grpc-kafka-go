package repo

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd
}

type Repo struct {
	redisClient *redis.Client
}

func New() *Repo {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // server address
		Password: "",               // no password
		DB:       0,                // default DB
	})

	// Ping the Redis server to check if it's reachable
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Print("[Redis] Failed to connect:", err)
		return nil
	}

	return &Repo{
		redisClient: client,
	}
}
