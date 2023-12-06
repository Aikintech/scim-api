package config

import (
	"os"

	"github.com/gofiber/storage/redis/v3"
)

var RedisStore *redis.Storage

func InitializeRedis() {
	store := redis.New(redis.Config{
		URL:   os.Getenv("REDIS_URL"),
		Reset: false,
	})

	RedisStore = store
}
