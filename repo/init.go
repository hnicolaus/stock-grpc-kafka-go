/*
	Hans Nicolaus
	29 Aug 2023
*/

package repo

import (
	"context"
	"fmt"
	"log"

	"bibit.id/challenge/model"
	"github.com/go-redis/redis/v8"
)

//go:generate mockgen -source=./init.go -destination=./_mock/stock_summary_mock.go -package=mock
type RedisClient interface {
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
	ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd
	ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd
}

type Repo struct {
	redisClient RedisClient
}

func New(cfg model.Config) *Repo {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Ping the Redis server to check if it's reachable
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Print("[Redis] Failed to connect:", err)
		return nil
	}

	log.Printf("[Redis] Serving on port %s", cfg.Redis.Port)

	return &Repo{
		redisClient: client,
	}
}
