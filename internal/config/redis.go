package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	rdb *redis.Client
)

func ConnectRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
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
