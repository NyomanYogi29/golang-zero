package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	RDB *redis.Client
)

func ConnectRedis() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
