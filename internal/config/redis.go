package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	rdb *redis.Client
)

func ConnectRedis() error {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		return err
	}

	rdb = client

	return nil
}

func GetRedisClient() *redis.Client {
	return rdb
}
